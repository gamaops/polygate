package main

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

// ConsumerRedisStack is a pointer to rotate Redis client distribution list
type ConsumerRedisStack struct {
	readGroupArgs redis.XReadGroupArgs
	retryKeys     []string
	retryArgs     []interface{}
	index         int
	client        *redis.Client
	xretrySha     string
	available     int64
	consuming     int32
	wakeUpCh      chan bool
}

func (c *ConsumerRedisStack) WakeUp() {
	c.wakeUpCh <- true
}

// Consumer is a group of Redis clients with specifications about the service being consumed
type Consumer struct {
	service *ConfigurationServiceExpose
	stacks  []*ConsumerRedisStack
}

var consumers = make(map[string]*Consumer)
var consumersStopped = false
var consumersStopWait = sync.WaitGroup{}

func acquireJobs(consumer *Consumer, stack *ConsumerRedisStack) {

	stack.retryKeys[3] = strconv.FormatInt(stack.available, 10)

	retryMessages, err := stack.client.EvalSha(
		stack.xretrySha,
		stack.retryKeys,
		stack.retryArgs...,
	).Result()

	if err != nil && err != redis.Nil {
		log.Fatalf("Error while executing XRETRY script: %v", err)
	}

	switch retryMessagesTyped := retryMessages.(type) {
	case []interface{}:
		retryMessagesCount := -int64(len(retryMessagesTyped))
		atomic.AddInt64(&stack.available, retryMessagesCount)
		for _, retryMessage := range retryMessagesTyped {
			go parseRetryItemToJob(consumer, stack, retryMessage.([]interface{}))
		}
		break
	}

	if atomic.CompareAndSwapInt64(&stack.available, 0, 0) {
		atomic.StoreInt32(&stack.consuming, 0)
		return
	}

	stack.readGroupArgs.Count = stack.available
	readMessages, err := stack.client.XReadGroup(&stack.readGroupArgs).Result()

	// redis.Nil is when reading exists by timeout
	if err != nil && err != redis.Nil {
		log.Fatalf("Error while consuming streams: %v", err)
	}

	for _, readMessage := range readMessages {
		streamMessagesCount := -int64(len(readMessage.Messages))
		atomic.AddInt64(&stack.available, streamMessagesCount)
		for i := range readMessage.Messages {
			message := &readMessage.Messages[i]
			go parseStreamItemToJob(consumer, stack, readMessage.Stream, message.ID, message.Values)
		}
	}

	if atomic.CompareAndSwapInt64(&stack.available, 0, 0) {
		atomic.StoreInt32(&stack.consuming, 0)
		return
	}

	atomic.StoreInt32(&stack.consuming, 0)
	stack.WakeUp()

}

func acquireJobsLoop(consumer *Consumer, stack *ConsumerRedisStack) {

	for range stack.wakeUpCh {
		if consumersStopped {
			consumersStopWait.Done()
			return
		}
		if atomic.SwapInt32(&stack.consuming, 1) == 1 {
			continue
		}
		go acquireJobs(consumer, stack)
	}

}

func ensureConsumerGroup(stream string, group string, service string) {

	for _, client := range redisClients[service] {
		_, err := client.XGroupCreateMkStream(stream, group, "$").Result()
		if err != nil {
			log.Warnf("Probably not critical, but create group/stream caused an error: %v", err)
		}

	}

}

func buildConsumerRedisStack(
	clients map[int]*redis.Client,
	readGroupArgs redis.XReadGroupArgs,
	service *ConfigurationServiceExpose,
) []*ConsumerRedisStack {

	xretryLoadScript, err := Asset("xretry.lua")

	if err != nil {
		log.Fatalf("Error while trying to load XRETRY script: %v", err)
	}

	clientsCount := len(clients)
	stacks := make([]*ConsumerRedisStack, clientsCount)
	concurrency := int64(int(service.Consumer.Concurrency) / clientsCount)

	deadline, err := time.ParseDuration(service.Consumer.Retry.Deadline)

	if err != nil {
		log.Fatalf("Invalid duration for retry deadline: %v", err)
	}

	retryArgs := make([]interface{}, len(readGroupArgs.Streams)/2)

	for index, stream := range readGroupArgs.Streams {
		if stream == ">" {
			break
		}
		retryArgs[index] = stream
	}

	for index, client := range clients {
		xretrySha, err := client.ScriptLoad(string(xretryLoadScript)).Result()
		if err != nil && err != redis.Nil {
			log.Fatalf("Error while loading script XRETRY: %v", err)
		}
		stacks[index] = &ConsumerRedisStack{
			readGroupArgs: redis.XReadGroupArgs{
				Group:    readGroupArgs.Group,
				Streams:  readGroupArgs.Streams,
				Consumer: readGroupArgs.Consumer,
				Block:    readGroupArgs.Block,
				NoAck:    false,
			},
			retryKeys: []string{
				readGroupArgs.Group,
				readGroupArgs.Consumer,
				strconv.FormatInt(service.Consumer.Retry.Limit, 10),
				"0",
				strconv.FormatInt(service.Consumer.Retry.PageSize, 10),
				strconv.FormatInt(deadline.Milliseconds(), 10),
			},
			retryArgs: retryArgs,
			client:    client,
			index:     index,
			xretrySha: xretrySha,
			available: concurrency,
			consuming: 0,
			wakeUpCh:  make(chan bool, 1),
		}
	}

	return stacks
}

func loadConsumers() {

	hostname, err := os.Hostname()

	if err != nil {
		log.Fatalf("Unable to get hostname: %v", err)
	}

	for i := range configuration.Protos.Services {

		service := &configuration.Protos.Services[i]

		startRedisClient(service.Service, 1)
		var streams []string
		var ids []string

		var group strings.Builder
		group.WriteString(configuration.Redis.Prefix)
		group.WriteRune(':')
		group.WriteString(service.Service)

		for _, method := range service.Methods {
			var stream strings.Builder
			stream.WriteString(configuration.Redis.Prefix)
			stream.WriteRune(':')
			stream.WriteString(method.Stream)
			ensureConsumerGroup(stream.String(), group.String(), service.Service)
			streams = append(streams, stream.String())
			ids = append(ids, ">")
		}

		streams = append(streams, ids...)

		var consumerID strings.Builder
		consumerID.WriteString(configuration.Redis.Prefix)
		consumerID.WriteString(":cons:")
		consumerID.WriteString(hostname)

		block, err := time.ParseDuration(service.Consumer.Block)

		if err != nil {
			log.Fatalf("Invalid duration: %v", err)
		}

		readGroupArgs := redis.XReadGroupArgs{
			Group:    group.String(),
			Streams:  streams,
			Consumer: consumerID.String(),
			Block:    block,
		}

		stacks := buildConsumerRedisStack(redisClients[service.Service], readGroupArgs, service)

		consumers[service.Service] = &Consumer{
			service: service,
			stacks:  stacks,
		}
	}

}

func startConsumers() {

	for i := range consumers {
		consumer := consumers[i]
		log.Debugf("Consuming service: %v", consumer.service.Service)
		for _, stack := range consumer.stacks {
			stack.WakeUp()
			consumersStopWait.Add(1)
			go acquireJobsLoop(consumer, stack)
		}
	}

}

func stopConsumersRedisConnections() {

	for i := range consumers {
		consumer := consumers[i]
		closeRedisClients(consumer.service.Service)
	}

}

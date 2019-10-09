package main

import (
	"fmt"
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
	previous      *ConsumerRedisStack
	xretrySha     string
	lock          sync.Mutex
	linked        bool
}

// Consumer is a group of Redis clients with specifications about the service being consumed
type Consumer struct {
	service        *ConfigurationServiceExpose
	maxCount       uint32
	pendingCount   int32
	availableCount uint32
	stack          *ConsumerRedisStack
	lock           sync.Mutex
}

// StreamItem represents a stream item for XRETRY command
type StreamItem struct {
	id     string
	stream string
	data   map[string]interface{}
}

var consumers = make(map[string]*Consumer)

func acquireJobs(consumer *Consumer, stack *ConsumerRedisStack) {

	stack.lock.Lock()

	stack.retryKeys[3] = strconv.FormatInt(stack.readGroupArgs.Count, 10)

	retryMessages, err := stack.client.EvalSha(
		stack.xretrySha,
		stack.retryKeys,
		stack.retryArgs...,
	).Result()

	if err != nil {
		log.Fatalf("Error while executing XRETRY script: %v", err)
	}

	// TODO: atomic add pendingCount
	switch retryMessagesTyped := retryMessages.(type) {
	case []interface{}:
		retryMessagesCount := int32(len(retryMessagesTyped))
		atomic.AddInt32(&consumer.pendingCount, retryMessagesCount)
		stack.readGroupArgs.Count -= int64(retryMessagesCount)
		for _, retryMessage := range retryMessagesTyped {
			// TODO: Handle retry jobs
			fmt.Println(retryMessage.(StreamItem))
		}
		break
	}

	if stack.readGroupArgs.Count == 0 {
		stack.lock.Unlock()
		return
	}

	if err != nil && err != redis.Nil {
		log.Fatalf("Error while consuming streams (retry): %v", err)
	}

	log.Infof("Retry: %v", retryMessages)

	log.Info("Consuming")

	readMessages, err := stack.client.XReadGroup(&stack.readGroupArgs).Result()

	// redis.Nil is when reading exists by timeout
	if err != nil && err != redis.Nil {
		log.Fatalf("Error while consuming streams: %v", err)
	}

	var messagesCount int32 = 0

	for _, readMessage := range readMessages {
		for _, message := range readMessage.Messages {
			job := parseStreamItemToJob(message.ID, message.Values)
			job.stack = stack
			job.consumer = consumer
			go clientJobHandlers[job.event.Service][job.event.Method](job)
		}
		streamMessagesCount := int32(len(readMessage.Messages))
		messagesCount += streamMessagesCount
		stack.readGroupArgs.Count -= int64(streamMessagesCount)
	}

	if messagesCount > 0 {
		atomic.AddInt32(&consumer.pendingCount, messagesCount)
	}

	if stack.readGroupArgs.Count > 0 {
		defer consumeService(consumer.service)
		defer linkRedisStack(consumer, stack)
	}

	stack.lock.Unlock()

}

func linkRedisStack(consumer *Consumer, stack *ConsumerRedisStack) {
	consumer.lock.Lock()

	defer consumer.lock.Unlock()

	if stack.linked {
		return
	}

	stack.linked = true
	availableCount := 1
	first := stack.previous

	for first != nil {
		availableCount++
		if first.previous == nil {
			first.previous = consumer.stack
			break
		}
		first = first.previous
	}

	consumer.stack = stack

	atomic.AddUint32(&consumer.availableCount, uint32(availableCount))

}

func consumeService(service *ConfigurationServiceExpose) {

	consumer := consumers[service.Service]

	consumer.lock.Lock()

	freeSlots := consumer.service.Consumer.Concurrency - uint32(consumer.pendingCount)

	if consumer.stack == nil || consumer.availableCount == 0 || freeSlots <= 0 {
		consumer.lock.Unlock()
		return
	}

	availableCount := atomic.SwapUint32(&consumer.availableCount, 0)
	stack := consumer.stack
	consumer.stack = nil

	consumer.lock.Unlock()

	countPerClient := freeSlots / availableCount
	if countPerClient > consumer.maxCount {
		countPerClient = consumer.maxCount
	}
	modulo := int32(countPerClient % 1)

	if countPerClient < 1 {
		countPerClient = 0
		modulo = int32(freeSlots)
	} else if modulo > 0 {
		modulo = int32(availableCount) * modulo
	}

	for stack != nil {
		cursor := stack
		cursor.linked = false
		count := countPerClient
		if atomic.AddInt32(&modulo, -1) > -1 {
			count++
		}
		if count == 0 {
			linkRedisStack(consumer, cursor)
			break
		}
		previous := cursor.previous
		cursor.previous = nil
		cursor.readGroupArgs.Count = int64(count)
		go acquireJobs(consumer, cursor)
		stack = previous
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
) *ConsumerRedisStack {

	xretryLoadScript, err := Asset("xretry.lua")

	if err != nil {
		log.Fatalf("Error while trying to load XRETRY script: %v", err)
	}

	var stack *ConsumerRedisStack

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
		next := &ConsumerRedisStack{
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
			previous:  stack,
			xretrySha: xretrySha,
			lock:      sync.Mutex{},
			linked:    true,
		}
		stack = next
	}

	return stack
}

func loadConsumers() {

	hostname, err := os.Hostname()

	if err != nil {
		log.Fatalf("Unable to get hostname: %v", err)
	}

	for _, service := range configuration.Protos.Services {

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

		stack := buildConsumerRedisStack(redisClients[service.Service], readGroupArgs, &service)
		availableCount := uint32(len(redisClients[service.Service]))

		consumers[service.Service] = &Consumer{
			service:        &service,
			maxCount:       service.Consumer.Concurrency / availableCount,
			pendingCount:   0,
			availableCount: availableCount,
			stack:          stack,
			lock:           sync.Mutex{},
		}
	}

}

func startConsumers() {

	for _, consumer := range consumers {
		go consumeService(consumer.service)
	}

}

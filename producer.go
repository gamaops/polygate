package main

import (
	"strings"

	polygate_data "polygate/polygate-data"

	b64 "encoding/base64"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

var jobAwaitChannels = make(map[string]chan *polygate_data.JobEvent)

func logEvent(job *Job, data []byte) {
	log.WithFields(map[string]interface{}{
		"source": "event",
		"data":   b64.URLEncoding.EncodeToString(data),
		"jobId":  job.event.Id,
	}).Info("New event received")
}

func sendJob(job *Job, method *ConfigurationMethodExpose) string {

	jobClient := routeRedisClient("job", []byte(job.event.Id))

	data, err := proto.Marshal(job.event)

	if err != nil {
		log.Fatalf("Error while encoding event to protobuf: %v", err)
	}

	logEvent(job, data)

	values := map[string]interface{}{
		"event": data,
	}

	var stream strings.Builder

	stream.WriteString(configuration.Redis.Prefix)
	stream.WriteRune(':')
	stream.WriteString(job.event.Stream)

	args := redis.XAddArgs{
		Stream: stream.String(),
		Values: values,
	}

	if method.Capped > 0 {
		args.MaxLenApprox = int64(method.Capped)
	}

	streamID, err := jobClient.XAdd(&args).Result()

	if err != nil {
		log.Fatalf("Error while sending job: %v", err)
	}

	return streamID

}

func sendJobAndAwait(job *Job, method *ConfigurationMethodExpose) *polygate_data.JobEvent {

	jobAwaitChannels[job.event.Id] = make(chan *polygate_data.JobEvent)
	go sendJob(job, method)
	event := <-jobAwaitChannels[job.event.Id]
	delete(jobAwaitChannels, job.event.Id)
	return event

}

func receiveMessagesFromPubSub(channel <-chan *redis.Message) {
	for message := range channel {
		event := &polygate_data.JobEvent{}
		err := proto.Unmarshal([]byte(message.Payload), event)
		if err != nil {
			log.Errorf("Received invalid event, unable to decode message: %v", err)
			continue
		}
		log.Debugf("Received message from channel: %v", event.Id)
		awaitChannel, ok := jobAwaitChannels[event.Id]
		if ok {
			awaitChannel <- event
			continue
		}
		log.Warnf("Unable to find a listener to event, it'll be discarded: %v", event.Id)
	}
}

func startProducerListener() {

	pubsubClients["producer"] = make(map[int]*PubSubClient, len(redisClients["producer"]))

	for index, client := range redisClients["producer"] {
		pubsub := client.Subscribe(instanceProducerID)
		pubsubClients["producer"][index] = &PubSubClient{
			pubsub: pubsub,
		}
		channel := pubsub.Channel()
		go receiveMessagesFromPubSub(channel)
	}

}

func closeProducerListener() {

	log.Warnf("Closing producer clients")

	for _, pubsubClient := range pubsubClients["producer"] {
		err := pubsubClient.pubsub.Unsubscribe(instanceProducerID)
		if err != nil {
			log.Errorf("Error while unsubscribing producer client: %v", err)
			continue
		}
		err = pubsubClient.pubsub.Close()
		if err != nil {
			log.Errorf("Error while closing producer client: %v", err)
		}
	}

	log.Warnf("Producer clients closed")

	delete(pubsubClients, "producer")

}

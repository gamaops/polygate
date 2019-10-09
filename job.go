package main

import (
	"errors"
	polygate_data "polygate/polygate-data"
	"sync/atomic"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
)

// ReceivedJob defines methods to received jobs to be completed
type ReceivedJob interface {
	Client() *redis.Client
	Ack() error
	Resolve() error
	Reject([]byte) error
}

// Job holds pointers and data to represent a task to Polygate
type Job struct {
	event    *polygate_data.JobEvent
	stack    *ConsumerRedisStack
	consumer *Consumer
}

// Client returns a job handler Redis client
func (j *Job) Client() *redis.Client {
	client := redisClients["job"][j.stack.index]
	return client
}

// Ack acknowledges job on stream
func (j *Job) Ack() error {
	if j.stack == nil {
		return errors.New("job: job is already resolved, you can't ack")
	}
	_, err := j.Client().XAck(
		j.event.Stream,
		j.stack.readGroupArgs.Group,
		j.event.StreamId,
	).Result()
	return err
}

// Resolve resolves job and send the feedback to pubsub channel
func (j *Job) Resolve() error {
	if j.stack == nil {
		return errors.New("job: job is already resolved, you can't resolve it again")
	}
	client := j.Client()
	stack := j.stack
	j.stack = nil

	j.event.ConsumerId = stack.readGroupArgs.Consumer
	j.event.Group = stack.readGroupArgs.Group
	j.event.Status = polygate_data.JobEvent_RESOLVED

	data, err := proto.Marshal(j.event)

	if err != nil {
		return err
	}

	_, err = client.Publish(
		j.event.ProducerId,
		data,
	).Result()

	atomic.AddInt32(&j.consumer.pendingCount, -1)

	return err
}

// Reject rejects job and send the feedback to pubsub channel
func (j *Job) Reject() error {
	if j.stack == nil {
		return errors.New("job: job is already resolved, you can't reject it")
	}
	client := j.Client()
	stack := j.stack
	j.stack = nil

	j.event.ConsumerId = stack.readGroupArgs.Consumer
	j.event.Group = stack.readGroupArgs.Group
	j.event.Status = polygate_data.JobEvent_REJECTED

	data, err := proto.Marshal(j.event)

	if err != nil {
		return err
	}

	_, err = client.Publish(
		j.event.ProducerId,
		data,
	).Result()

	atomic.AddInt32(&j.consumer.pendingCount, -1)

	return err
}

// Reset function to protobuf parser
func (j *Job) Reset() { *j = Job{} }

// String function to protobuf parser
func (j *Job) String() string { return proto.CompactTextString(j) }

// ProtoMessage function to protobuf parser
func (*Job) ProtoMessage() {}

// XXX_Unmarshal function to protobuf parser
func (j *Job) XXX_Unmarshal(b []byte) error {
	j.event = &polygate_data.JobEvent{
		Payload: b,
	}
	return nil
}

// XXX_Marshal function to protobuf parser
func (j *Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(j.event.Payload, nil, deterministic)
}

// XXX_Merge function to protobuf parser
func (j *Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(j, src)
}

// XXX_Size function to protobuf parser
func (j *Job) XXX_Size() int {
	return xxx_messageInfo_Message.Size(j)
}

// XXX_DiscardUnknown function to protobuf parser
func (j *Job) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(j)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

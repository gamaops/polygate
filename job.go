package main

import (
	"errors"
	polygate_data "polygate/polygate-data"
	"sync/atomic"

	"google.golang.org/grpc"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
)

// ReceivedJob defines methods to received jobs to be completed
type ReceivedJob interface {
	Client() *redis.Client
	Ack() error
	Resolve() error
	Reject() error
	Finish() error
}

// Job holds pointers and data to represent a task to Polygate
type Job struct {
	event     *polygate_data.JobEvent
	stack     *ConsumerRedisStack
	consumer  *Consumer
	rawStream string
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
		j.rawStream,
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

	atomic.AddInt64(&stack.available, 1)

	if err != nil {
		stack.WakeUp()
		return err
	}

	_, err = client.Publish(
		j.event.ProducerId,
		data,
	).Result()

	if err == nil {
		stack.WakeUp()
	}

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

	atomic.AddInt64(&stack.available, 1)

	if err != nil {
		stack.WakeUp()
		return err
	}

	_, err = client.Publish(
		j.event.ProducerId,
		data,
	).Result()

	if err == nil {
		stack.WakeUp()
	}

	return err
}

// Finish finishes the job
func (j *Job) Finish() error {
	if j.stack == nil {
		return errors.New("job: job is already resolved, you can't reject it")
	}

	stack := j.stack
	j.stack = nil

	atomic.AddInt64(&stack.available, 1)

	stack.WakeUp()

	return nil
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

type PolygateClientStreamServer interface {
	SendAndClose(*Job) error
	Recv() (*Job, error)
	grpc.ServerStream
}

type polygateClientStream struct {
	grpc.ServerStream
}

func (x *polygateClientStream) SendAndClose(j *Job) error {
	return x.ServerStream.SendMsg(j)
}

func (x *polygateClientStream) Recv() (*Job, error) {
	j := new(Job)
	if err := x.ServerStream.RecvMsg(j); err != nil {
		return nil, err
	}
	return j, nil
}

type PolygateClientStreamClient interface {
	SendAndClose(*Job) error
	Recv() (*Job, error)
	grpc.ClientStream
}

type polygateClientStreamClient struct {
	grpc.ClientStream
}

func (x *polygateClientStreamClient) Send(j *Job) error {
	return x.ClientStream.SendMsg(j)
}

func (x *polygateClientStreamClient) CloseAndRecv() (*Job, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	j := new(Job)
	if err := x.ClientStream.RecvMsg(j); err != nil {
		return nil, err
	}
	return j, nil
}

var emptyJob = &Job{
	event: &polygate_data.JobEvent{
		Payload:  make([]byte, 0),
		Metadata: make([]*polygate_data.MetadataItem, 0),
	},
}

package main

import (
	"context"
	polygate_data "polygate/polygate-data"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type MethodClientStream struct {
	stream *polygateClientStreamClient
	ctx    context.Context
	timer  *ResetableTimer
}

type UpstreamMethod struct {
	route            string
	clientStreamDesc *grpc.StreamDesc
	clientStreams    sync.Map
	lock             sync.Mutex
}

type ClientUpstream struct {
	client  *grpc.ClientConn
	service *ConfigurationServiceExpose
	methods map[string]*UpstreamMethod
}

var clientJobHandlers map[string]map[string]func(*Job)
var clientUpstreams map[string]*ClientUpstream

func parseStreamItemToJob(consumer *Consumer, stack *ConsumerRedisStack, rawStream string, id string, data map[string]interface{}) {
	event := &polygate_data.JobEvent{}
	err := proto.Unmarshal([]byte(data["event"].(string)), event)
	if err != nil {
		log.Fatalf("Unable to parse data from job %v: $v", id, err)
	}
	event.StreamId = id
	job := &Job{
		event:     event,
		rawStream: rawStream,
		stack:     stack,
		consumer:  consumer,
	}
	clientJobHandlers[job.event.Service][job.event.Method](job)
}

func parseRetryItemToJob(consumer *Consumer, stack *ConsumerRedisStack, item []interface{}) {
	data := item[5].([]interface{})[1].(string)
	event := &polygate_data.JobEvent{}
	err := proto.Unmarshal([]byte(data), event)
	if err != nil {
		log.Fatalf("Unable to parse data from job (streamId %v): $v", item[0].(string), err)
	}
	event.StreamId = item[0].(string)
	log.WithFields(map[string]interface{}{
		"streamId":    event.StreamId,
		"jobId":       event.Id,
		"timeWaiting": item[2],
		"retries":     item[3],
	}).Warnf("Retrying job")
	job := &Job{
		event:     event,
		rawStream: item[4].(string),
		stack:     stack,
		consumer:  consumer,
	}
	clientJobHandlers[job.event.Service][job.event.Method](job)
}

func loadClientConn(service *ConfigurationServiceExpose) *grpc.ClientConn {

	var location strings.Builder
	location.WriteString(service.Client.Address)
	location.WriteRune(':')
	location.WriteString(strconv.FormatUint(uint64(service.Client.Port), 10))

	// TODO: Maybe in the future implement secure context?
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.ForceCodecCallOption{
				Codec: Codec{},
			},
		),
	}

	conn, err := grpc.Dial(location.String(), opts...)

	if err != nil {
		log.Fatalf("Fail to dial: %v", err)
	}

	return conn

}

func ensureJobClientStream(upstream *ClientUpstream, method *UpstreamMethod, md *metadata.MD, timeout time.Duration) (*MethodClientStream, error) {

	callID := md.Get("callId")[0]

	methodClientStream, ok := method.clientStreams.Load(callID)

	if ok {
		return methodClientStream.(*MethodClientStream), nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	stream, err := upstream.client.NewStream(ctx, method.clientStreamDesc, method.route)

	if err != nil {
		cancel()
		return nil, err
	}

	clientStream := &polygateClientStreamClient{stream}

	newClientStream := &MethodClientStream{
		stream: clientStream,
		ctx:    ctx,
		timer:  NewResetableTimer(timeout, 1),
	}

	go func() {
		status := <-newClientStream.timer.Status
		if status == RTTimeout {
			method.clientStreams.Delete(callID)
			log.Debugf("Stream timeout while waiting for next: %v", callID)

			go func() {
				_, err := newClientStream.stream.CloseAndRecv()
				newClientStream.timer.Status <- RTCancel
				cancel()
				if err != nil {
					log.WithFields(map[string]interface{}{
						"callId": callID,
					}).Warnf("Client stream rejected: %v", err)
				} else {
					log.WithFields(map[string]interface{}{
						"callId": callID,
					}).Info("Client stream resolved")
				}
			}()
			go newClientStream.timer.Start()

			switch <-newClientStream.timer.Status {
			case RTTimeout:
				log.WithFields(map[string]interface{}{
					"callId": callID,
				}).Warn("Stream cancelled")
				cancel()
			}

		}
	}()

	method.clientStreams.Store(callID, newClientStream)

	return newClientStream, nil

}

func loadClientJobHandlers() {

	clientJobHandlers = make(map[string]map[string]func(*Job), len(configuration.Protos.Services))
	clientUpstreams = make(map[string]*ClientUpstream, len(configuration.Protos.Services))
	for i := range configuration.Protos.Services {

		service := &configuration.Protos.Services[i]

		clientUpstream := &ClientUpstream{
			client:  loadClientConn(service),
			service: service,
			methods: make(map[string]*UpstreamMethod, len(service.Methods)),
		}

		clientUpstreams[service.Service] = clientUpstream

		methodsHandlers := make(map[string]func(*Job), len(service.Methods))
		for t := range service.Methods {

			method := &service.Methods[t]

			var methodRoute strings.Builder
			methodRoute.WriteRune('/')
			methodRoute.WriteString(service.Service)
			methodRoute.WriteRune('/')
			methodRoute.WriteString(method.Name)

			upstreamMethod := &UpstreamMethod{
				route:            methodRoute.String(),
				clientStreamDesc: nil,
				clientStreams:    sync.Map{},
				lock:             sync.Mutex{},
			}

			clientUpstream.methods[method.Name] = upstreamMethod

			if method.Pattern == "queue" {
				methodsHandlers[method.Name] = func(job *Job) {

					err := job.Ack()
					log.Debugf("Acknowledge job: %v", job.event.Id)
					if err != nil {
						log.Fatalf("Error on job acknowledgement: %v", err)
					}

					out := new(Job)
					upstream := clientUpstreams[job.event.Service]

					ctx := context.Background()
					md := metadataFromJobEvent(job.event)

					ctx = metadata.NewOutgoingContext(ctx, md)

					var header metadata.MD

					err = upstream.client.Invoke(ctx, upstreamMethod.route, job, out, grpc.Header(&header))

					addMetadataToJobEvent(header, job.event)

					if err != nil {

						addErrorToJobEvent(err, job.event)
						job.event.Payload = make([]byte, 0)
						err = job.Reject()

						if err != nil {
							log.Fatalf("Error publishing rejection: %v", err)
						}
						return
					}

					job.event.Payload = out.event.Payload
					err = job.Resolve()

					if err != nil {
						log.Fatalf("Error publishing resolution: %v", err)
					}
				}
			} else if method.Pattern == "fireAndForget" {
				upstreamMethod.clientStreamDesc = &grpc.StreamDesc{
					StreamName:    method.Name,
					ClientStreams: true,
				}
				timeoutWaitForNext, err := time.ParseDuration(method.TimeoutWaitForNext)
				if err != nil {
					log.Fatalf("Invalid duration for timeoutWaitForNext: %v", err)
				}
				methodsHandlers[method.Name] = func(job *Job) {

					md := metadataFromJobEvent(job.event)

					methodClientStream, err := ensureJobClientStream(
						clientUpstream,
						upstreamMethod,
						&md,
						timeoutWaitForNext,
					)

					callID := md.Get("callId")[0]

					if err != nil {
						log.WithFields(map[string]interface{}{
							"callId": callID,
							"jobId":  job.event.Id,
						}).Warnf("Fire and forget job rejection while acquiring a new stream: %v", err)
						addErrorToJobEvent(err, job.event)
						job.event.Payload = make([]byte, 0)
						err = job.Reject()
						if err != nil {
							log.Fatalf("Error publishing rejection: %v", err)
						}
						return
					}

					err = methodClientStream.stream.Send(job)

					if err != nil {
						log.WithFields(map[string]interface{}{
							"callId": callID,
							"jobId":  job.event.Id,
						}).Warnf("Fire and forget job rejection while sending job to upstream: %v", err)
						// addErrorToJobEvent(err, job.event)
						// job.event.Payload = make([]byte, 0)
						// err = job.Reject()
						// if err != nil {
						// 	log.Fatalf("Error publishing rejection: %v", err)
						// }
						return
					}

					methodClientStream.timer.Reset()

					err = job.Ack()
					if err != nil {
						log.Fatalf("Error on job acknowledgement: %v", err)
					}

					err = job.Finish()
					if err != nil {
						log.Fatalf("Error publishing resolution: %v", err)
					}

				}
			}
		}

		clientJobHandlers[service.Service] = methodsHandlers
	}

}

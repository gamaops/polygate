package main

import (
	"context"
	polygate_data "polygate/polygate-data"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientUpstream struct {
	client       *grpc.ClientConn
	service      *ConfigurationServiceExpose
	methodsRoute map[string]string
}

var clientJobHandlers map[string]map[string]func(*Job)
var clientUpstreams map[string]*ClientUpstream

func parseStreamItemToJob(id string, data map[string]interface{}) *Job {
	event := &polygate_data.JobEvent{}
	err := proto.Unmarshal([]byte(data["event"].(string)), event)
	if err != nil {
		log.Fatalf("Unable to parse data from job %v: $v", id, err)
	}
	event.StreamId = id
	return &Job{
		event: event,
	}
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

func loadClientJobHandlers() {

	clientJobHandlers = make(map[string]map[string]func(*Job), len(configuration.Protos.Services))
	clientUpstreams = make(map[string]*ClientUpstream, len(configuration.Protos.Services))
	for _, service := range configuration.Protos.Services {

		clientUpstreams[service.Service] = &ClientUpstream{
			client:       loadClientConn(&service),
			service:      &service,
			methodsRoute: make(map[string]string, len(service.Methods)),
		}

		methodsHandlers := make(map[string]func(*Job), len(service.Methods))
		for _, method := range service.Methods {

			var methodRoute strings.Builder
			methodRoute.WriteRune('/')
			methodRoute.WriteString(service.Service)
			methodRoute.WriteRune('/')
			methodRoute.WriteString(method.Name)
			clientUpstreams[service.Service].methodsRoute[method.Name] = methodRoute.String()

			if method.Pattern == "queue" {
				methodsHandlers[method.Name] = func(job *Job) {

					job.Ack()
					out := new(Job)
					upstream := clientUpstreams[job.event.Service]

					log.Infof("Job received: %v", job)

					ctx := context.Background()
					md := metadataFromJobEvent(job.event)

					ctx = metadata.NewOutgoingContext(ctx, md)

					var header metadata.MD

					err := upstream.client.Invoke(ctx, upstream.methodsRoute[method.Name], job, out, grpc.Header(&header))

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
			}
		}

		clientJobHandlers[service.Service] = methodsHandlers
	}

}

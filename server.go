package main

import (
	"context"
	"fmt"
	"io"
	"net"

	polygate_data "polygate/polygate-data"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func createServiceServer(server *grpc.Server, expose *ConfigurationServiceExpose) {

	var methods []grpc.MethodDesc
	var streams []grpc.StreamDesc
	var handlers map[string]interface{}

	for i := range expose.Methods {
		method := &expose.Methods[i]
		mJobCount := producerJobCount.WithLabelValues(expose.Service, method.Name, method.Stream)
		mJobExecutionSeconds := producerJobExecutionSeconds.WithLabelValues(expose.Service, method.Name, method.Stream)
		mJobPayloadBytes := producerJobPayloadBytes.WithLabelValues(expose.Service, method.Name, method.Stream)
		mJobEventBytes := producerJobEventBytes.WithLabelValues(expose.Service, method.Name, method.Stream)
		if method.Pattern == "queue" {
			mFailedJobCount := producerFailedJobCount.WithLabelValues(expose.Service, method.Name, method.Stream)
			handler := func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {

				mJobCount.Inc()
				timer := prometheus.NewTimer(mJobExecutionSeconds)
				defer timer.ObserveDuration()

				in := new(Job)
				err := dec(in)

				if err != nil {
					mFailedJobCount.Inc()
					return nil, err
				}

				in.event.Id = xid.New().String()
				in.event.Service = expose.Service
				in.event.Method = method.Name
				in.event.Stream = method.Stream
				in.event.ProducerId = instanceProducerID
				in.event.Status = polygate_data.JobEvent_AWAITING

				md, ok := metadata.FromIncomingContext(ctx)

				if ok {
					addMetadataToJobEvent(md, in.event)
				}

				outEvent := sendJobAndAwait(in, method, mJobPayloadBytes, mJobEventBytes)
				md = metadataFromJobEvent(outEvent)
				grpc.SendHeader(ctx, md)

				if outEvent.Status == polygate_data.JobEvent_REJECTED {
					st := statusFromJobEvent(outEvent)
					mFailedJobCount.Inc()
					return nil, st.Err()
				}

				return outEvent.Payload, nil
			}
			methods = append(methods, grpc.MethodDesc{
				MethodName: method.Name,
				Handler:    handler,
			})
		} else if method.Pattern == "fireAndForget" {
			mClientStreamsCount := producerClientStreamsCount.WithLabelValues(expose.Service, method.Name, method.Stream)
			handler := func(srv interface{}, stream grpc.ServerStream) error {

				mClientStreamsCount.Inc()
				timer := prometheus.NewTimer(mJobExecutionSeconds)
				defer timer.ObserveDuration()

				clientStream := &polygateClientStream{stream}
				callID := xid.New().String()

				log.WithFields(map[string]interface{}{
					"type":   "clientStream",
					"callId": callID,
				}).Info("New client stream call")

				for {

					mJobCount.Inc()

					in, err := clientStream.Recv()
					if err == io.EOF {
						// TODO: Send as metadata callId
						mClientStreamsCount.Dec()
						return clientStream.SendAndClose(emptyJob)
					}
					if err != nil {
						mClientStreamsCount.Dec()
						return err
					}

					in.event.Id = xid.New().String()
					in.event.Service = expose.Service
					in.event.Method = method.Name
					in.event.Stream = method.Stream
					in.event.ProducerId = instanceProducerID
					in.event.Status = polygate_data.JobEvent_AWAITING

					md, ok := metadata.FromIncomingContext(stream.Context())
					md.Set("callId", callID)

					if ok {
						addMetadataToJobEvent(md, in.event)
					}

					log.WithFields(map[string]interface{}{
						"type":   "clientStream",
						"callId": callID,
						"jobId":  in.event.Id,
					}).Info("New job on client stream call")

					go sendJob(in, method, mJobPayloadBytes, mJobEventBytes)

				}
			}
			streams = append(streams, grpc.StreamDesc{
				StreamName:    method.Name,
				ClientStreams: true,
				Handler:       handler,
			})
		} else {
			log.Fatalf("Unknown pattern for service %v: %v", expose.Service, method.Pattern)
		}
	}

	serviceDescription := grpc.ServiceDesc{
		ServiceName: expose.Service,
		HandlerType: handlers,
		Methods:     methods,
		Streams:     streams,
	}

	server.RegisterService(&serviceDescription, handlers)
}

func createServer() *grpc.Server {

	options := []grpc.ServerOption{
		grpc.CustomCodec(Codec{}),
		grpc.MaxHeaderListSize(configuration.Server.MaxHeaderListSize),
	}

	server := grpc.NewServer(options...)

	for i := range configuration.Protos.Services {
		expose := &configuration.Protos.Services[i]
		createServiceServer(server, expose)
	}

	location := fmt.Sprintf("localhost:%d", configuration.Server.Port)
	listener, err := net.Listen("tcp", location)

	if err != nil {
		log.Fatalf("Failed to start server listener: %v", err)
	}

	go func() {
		err := server.Serve(listener)
		if err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	log.Infof("gRPC server is listening on: %v", location)

	producerReady.Set(1)

	return server

}

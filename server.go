package main

import (
	"context"
	"fmt"
	"net"

	polygate_data "polygate/polygate-data"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func createServiceServer(server *grpc.Server, expose *ConfigurationServiceExpose) {

	var methods []grpc.MethodDesc
	var streams []grpc.StreamDesc
	var handlers map[string]interface{}

	for _, method := range expose.Methods {
		if method.Pattern == "queue" { // Unary Call Handler
			handler := func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {

				in := new(Job)
				err := dec(in)

				if err != nil {
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

				outEvent := sendJobAndAwait(in, &method)
				md = metadataFromJobEvent(outEvent)
				grpc.SendHeader(ctx, md)

				if outEvent.Status == polygate_data.JobEvent_REJECTED {
					st := statusFromJobEvent(outEvent)
					return nil, st.Err()
				}

				return outEvent.Payload, nil
			}
			methods = append(methods, grpc.MethodDesc{
				MethodName: method.Name,
				Handler:    handler,
			})
		} else if method.Pattern == "fireAndForget" { // TODO: Client Stream Handler
			streams = append(streams, grpc.StreamDesc{
				StreamName:    method.Name,
				ClientStreams: true,
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

	for _, expose := range configuration.Protos.Services {
		createServiceServer(server, &expose)
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

	return server

}

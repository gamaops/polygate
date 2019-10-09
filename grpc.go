package main

import (
	polygate_data "polygate/polygate-data"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func addMetadataToJobEvent(md metadata.MD, event *polygate_data.JobEvent) {

	event.Metadata = make([]*polygate_data.MetadataItem, len(md))
	index := 0
	for key, value := range md {
		event.Metadata[index] = &polygate_data.MetadataItem{
			Key:    key,
			Values: value,
		}
		index++
	}

}

func metadataFromJobEvent(event *polygate_data.JobEvent) metadata.MD {

	md := make(metadata.MD, len(event.Metadata))

	for _, value := range event.Metadata {
		md.Set(value.Key, value.Values...)
	}

	return md

}

func addErrorToJobEvent(err error, event *polygate_data.JobEvent) {
	s := status.Convert(err)

	event.Error = &polygate_data.JobError{
		Code:    uint32(s.Code()),
		Message: s.Message(),
	}

}

func statusFromJobEvent(event *polygate_data.JobEvent) *status.Status {
	return status.New(codes.Code(event.Error.Code), event.Error.Message)
}

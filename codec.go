package main

import (
	"github.com/golang/protobuf/proto"
)

type Codec struct{}

func (Codec) Marshal(v interface{}) ([]byte, error) {
	switch vTyped := v.(type) {
	case *Job:
		return vTyped.event.Payload, nil
	case []uint8:
		return []byte(vTyped), nil
	}
	return v.([]byte), nil
}

func (Codec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (Codec) String() string {
	return "no-op-marshal codec"
}

func (Codec) Name() string {
	return "no-op-marshal codec"
}

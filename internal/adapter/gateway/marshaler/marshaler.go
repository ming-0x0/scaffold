package marshaler

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type MarshalerInterface interface {
	NewNilMarshaler() runtime.Marshaler
}

type Marshaler struct{}

func New() MarshalerInterface {
	return &Marshaler{}
}

type NilMarshaler struct {
	runtime.Marshaler
}

func (m *NilMarshaler) Marshal(v any) ([]byte, error) {
	return nil, nil
}

func (m *NilMarshaler) Unmarshal(data []byte, v any) error {
	return nil
}

func (m *Marshaler) NewNilMarshaler() runtime.Marshaler {
	return &NilMarshaler{
		Marshaler: &runtime.JSONPb{},
	}
}

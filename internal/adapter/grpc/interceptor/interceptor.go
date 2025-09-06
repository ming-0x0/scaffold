package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type InterceptorKey string

const (
	requestIDKey InterceptorKey = "request_id"
)

type InterceptorInterface interface {
	InterceptContext(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error)
}

type Interceptor struct{}

func New() InterceptorInterface {
	return &Interceptor{}
}

func (i *Interceptor) InterceptContext(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if requestIDSlice := md.Get(string(requestIDKey)); len(requestIDSlice) > 0 {
			requestID := requestIDSlice[0]
			ctx = context.WithValue(ctx, requestIDKey, requestID)
		}
	}

	return handler(ctx, req)
}

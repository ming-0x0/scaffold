package interceptor

import (
	"context"
	"errors"
	"runtime/debug"

	"github.com/ming-0x0/scaffold/internal/shared/domainerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type InterceptorKey string

const (
	requestIDKey InterceptorKey = "request_id"
)

type InterceptorInterface interface {
	InterceptContext(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error)
	InterceptPanic(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error)
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

func (i *Interceptor) InterceptPanic(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = domainerror.Wrap(
				domainerror.Internal,
				errors.New(string(debug.Stack())),
			).GRPCStatus()
		}
	}()

	return handler(ctx, req)
}

package annotator

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

type AnnotatorKey string

const (
	requestIDKey AnnotatorKey = "request_id"
)

const requestIDHeader = "X-Request-ID"

type AnnotatorInterface interface {
	AnnotateMetadata(ctx context.Context, req *http.Request) metadata.MD
}

type Annotator struct{}

func New() AnnotatorInterface {
	return &Annotator{}
}

func (a *Annotator) AnnotateMetadata(ctx context.Context, req *http.Request) metadata.MD {
	md := metadata.MD{}

	if requestID := req.Header.Get(requestIDHeader); requestID != "" {
		md.Set(string(requestIDKey), requestID)
	}

	return md
}

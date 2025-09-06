package responder

import (
	"context"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

type ErrorWrapper struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

type ErrorResponderInterface interface {
	ErrorModifier(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error)
}

type ErrorResponder struct{}

func NewErrorResponder() *ErrorResponder {
	return &ErrorResponder{}
}

func (e *ErrorResponder) ErrorModifier(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	st := status.Convert(err)
	httpCode := grpcCodeToHTTPCode(st.Code())

	wrapped := ErrorWrapper{
		Code:    int(st.Code()),
		Message: st.Message(),
	}

	if os.Getenv("ENV") == "dev" {
		wrapped.Error = st.Details()
	}

	if err := writeJSONResponse(w, httpCode, wrapped); err != nil {
		return
	}
}

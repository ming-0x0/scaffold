package responder

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	common "github.com/ming-0x0/scaffold/api/gen/go/common"
	"google.golang.org/grpc/status"
)

type ErrorResponderInterface interface {
	RespondError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error)
}

type ErrorResponder struct{}

func NewErrorResponder() *ErrorResponder {
	return &ErrorResponder{}
}

func (e *ErrorResponder) RespondError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	st := status.Convert(err)
	httpCode := grpcCodeToHTTPCode(st.Code())

	wrapped := common.ErrorWrapper{
		Code:    int32(st.Code()),
		Message: st.Message(),
	}

	if os.Getenv("ENV") == "dev" {
		wrapped.Error = &common.ErrorDetails{
			Details: fmt.Sprintf("%v", st.Details()),
		}
	}

	if err := writeJSONResponse(w, httpCode, &wrapped); err != nil {
		return
	}
}

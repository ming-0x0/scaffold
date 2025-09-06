package responder

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

type ResponderInterface interface {
	ErrorResponderInterface
	ForwardResponderInterface
}

type Responder struct {
	errorResponder   *ErrorResponder
	forwardResponder *ForwardResponder
}

func New() *Responder {
	return &Responder{
		errorResponder:   NewErrorResponder(),
		forwardResponder: NewForwardResponder(),
	}
}

func (res *Responder) RespondError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	res.errorResponder.RespondError(ctx, mux, marshaler, w, r, err)
}

func (res *Responder) Respond(ctx context.Context, w http.ResponseWriter, msg proto.Message) error {
	return res.forwardResponder.Respond(ctx, w, msg)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func grpcCodeToHTTPCode(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusOK
	}
}

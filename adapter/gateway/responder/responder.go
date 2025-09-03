package responder

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ming-0x0/scaffold/api/gen/go/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type ResponderInterface interface {
	ErrorModifier(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error)
	ResponseModifier(ctx context.Context, w http.ResponseWriter, msg proto.Message) error
}

type Responder struct{}

func NewResponder() *Responder {
	return &Responder{}
}

type ErrorWrapper struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

type ResponseWrapper struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
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

func (res *Responder) ErrorModifier(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	st := status.Convert(err)
	httpCode := grpcCodeToHTTPCode(st.Code())

	wrapped := ErrorWrapper{
		Code:    int(st.Code()),
		Message: st.Message(),
	}

	if os.Getenv("ENV") == "dev" {
		wrapped.Error = st.Details()
	}

	err = writeJSONResponse(w, httpCode, wrapped)
	if err != nil {
		return
	}
}

func getResponseMessage(msg proto.Message) string {
	if msg == nil {
		return "OK"
	}

	// Get the message descriptor using ProtoReflect
	rf := msg.ProtoReflect()
	if !rf.IsValid() {
		return "OK"
	}

	// Get the descriptor
	desc := rf.Descriptor()
	if desc == nil {
		return "OK"
	}

	// Get the options
	opts := desc.Options()
	if opts == nil {
		return "OK"
	}

	// Get the response message option
	if ext, ok := proto.GetExtension(opts, common.E_ResponseMessage).(string); ok && ext != "" {
		return ext
	}

	return "OK"
}

func (res *Responder) ResponseModifier(ctx context.Context, w http.ResponseWriter, msg proto.Message) error {
	message := getResponseMessage(msg)

	wrapped := ResponseWrapper{
		Code:    0,
		Message: message,
		Data:    msg,
	}
	return writeJSONResponse(w, http.StatusOK, wrapped)
}

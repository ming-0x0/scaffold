package responder

import (
	"context"
	"net/http"

	"github.com/ming-0x0/scaffold/api/gen/go/common"
	"google.golang.org/protobuf/proto"
)

type ForwardResponseWrapper struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ForwardResponderInterface interface {
	Respond(ctx context.Context, w http.ResponseWriter, msg proto.Message) error
}

type ForwardResponder struct{}

func NewForwardResponder() *ForwardResponder {
	return &ForwardResponder{}
}

func getResponseMessage(msg proto.Message) string {
	if msg == nil {
		return "OK"
	}

	rf := msg.ProtoReflect()
	if !rf.IsValid() {
		return "OK"
	}

	desc := rf.Descriptor()
	if desc == nil {
		return "OK"
	}

	opts := desc.Options()
	if opts == nil {
		return "OK"
	}

	if ext, ok := proto.GetExtension(opts, common.E_ResponseMessage).(string); ok && ext != "" {
		return ext
	}

	return "OK"
}

func (f *ForwardResponder) Respond(ctx context.Context, w http.ResponseWriter, msg proto.Message) error {
	message := getResponseMessage(msg)

	wrapped := ForwardResponseWrapper{
		Code:    0,
		Message: message,
		Data:    msg,
	}

	return writeJSONResponse(w, http.StatusOK, wrapped)
}

package handler

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

const (
	RequestIDField = "request_id"
)

type RequestIDHandler struct {
	Handler slog.Handler
}

func (h *RequestIDHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if requestID, ok := ctx.Value(RequestIDField).(string); ok {
			r.AddAttrs(slog.String(RequestIDField, requestID))
		} else {
			r.AddAttrs(slog.String(RequestIDField, uuid.New().String()))
		}
	}
	return h.Handler.Handle(ctx, r)
}

func (h *RequestIDHandler) Enabled(ctx context.Context, r slog.Level) bool {
	return h.Handler.Enabled(ctx, r)
}

func (h *RequestIDHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return &RequestIDHandler{Handler: h.Handler}
}

func (h *RequestIDHandler) WithGroup(_ string) slog.Handler {
	return &RequestIDHandler{Handler: h.Handler}
}

func WithRequestIDHandler(handler slog.Handler) slog.Handler {
	return &RequestIDHandler{Handler: handler}
}

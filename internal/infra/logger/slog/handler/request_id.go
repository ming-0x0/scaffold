package handler

import (
	"context"
	"log/slog"
)

const (
	RequestIDField = "request_id"
)

type RequestIDHandler struct {
	handler slog.Handler
}

func (h *RequestIDHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if requestID, ok := ctx.Value(RequestIDField).(string); ok {
			r.AddAttrs(slog.String(RequestIDField, requestID))
		}
	}
	return h.handler.Handle(ctx, r)
}

func (h *RequestIDHandler) Enabled(ctx context.Context, r slog.Level) bool {
	return h.handler.Enabled(ctx, r)
}

func (h *RequestIDHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &RequestIDHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *RequestIDHandler) WithGroup(name string) slog.Handler {
	return &RequestIDHandler{handler: h.handler.WithGroup(name)}
}

func WithRequestID(handler slog.Handler) slog.Handler {
	return &RequestIDHandler{handler: handler}
}

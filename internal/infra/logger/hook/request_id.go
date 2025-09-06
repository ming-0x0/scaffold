package hook

import (
	"github.com/sirupsen/logrus"
)

const (
	RequestIDField = "request_id"
)

type RequestIDHook struct{}

func (hook *RequestIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *RequestIDHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	if requestID, ok := ctx.Value(RequestIDField).(string); ok {
		entry.Data[RequestIDField] = requestID
	}

	return nil
}

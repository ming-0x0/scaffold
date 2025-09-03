package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type LogLevel string

const (
	Panic LogLevel = "panic"
	Fatal LogLevel = "fatal"
	Error LogLevel = "error"
	Warn  LogLevel = "warn"
	Info  LogLevel = "info"
	Debug LogLevel = "debug"
	Trace LogLevel = "trace"
)

func (l LogLevel) ToLogrusLevel() logrus.Level {
	switch l {
	case Panic:
		return logrus.PanicLevel
	case Fatal:
		return logrus.FatalLevel
	case Error:
		return logrus.ErrorLevel
	case Warn:
		return logrus.WarnLevel
	case Info:
		return logrus.InfoLevel
	case Debug:
		return logrus.DebugLevel
	case Trace:
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}

func New() *logrus.Logger {
	logger := logrus.New()

	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     true,
	})
	logger.SetLevel(LogLevel(os.Getenv("LOG_LEVEL")).ToLogrusLevel())

	return logger
}

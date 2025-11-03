package sloglogger

import (
	"log/slog"

	"github.com/ming-0x0/scaffold/internal/infra/logger/slog/handler"
)

const (
	Panic = slog.Level(16)
	Fatal = slog.Level(12)
	Error = slog.Level(8)
	Warn  = slog.Level(4)
	Info  = slog.Level(0)
	Debug = slog.Level(-4)
	Trace = slog.Level(-8)
)

var levelNames = map[slog.Leveler]string{
	Panic: "panic",
	Fatal: "fatal",
	Error: "error",
	Warn:  "warn",
	Info:  "info",
	Debug: "debug",
	Trace: "trace",
}

func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, a.Value.Time().Format("15:04:05 02/01/2006"))
			}
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := levelNames[level]
				if !exists {
					levelLabel = level.String()
				}
				return slog.String(slog.LevelKey, levelLabel)
			}
			return a
		},
	}

	logger := handler.WithRequestIDHandler(NewJSONHandler(opts))

	return slog.New(logger)
}

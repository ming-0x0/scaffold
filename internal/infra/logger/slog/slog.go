package sloglogger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/ming-0x0/scaffold/internal/infra/logger/slog/handler"
	"github.com/ming-0x0/scaffold/pkg/env"
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

func getLogLevel(logLevel string) slog.Level {
	switch logLevel {
	case "panic":
		return Panic
	case "fatal":
		return Fatal
	case "error":
		return Error
	case "warn":
		return Warn
	case "info":
		return Info
	case "debug":
		return Debug
	case "trace":
		return Trace
	default:
		return Info
	}
}

type Logger struct {
	*slog.Logger
}

func New() *Logger {
	opts := &slog.HandlerOptions{
		Level:     getLogLevel(env.GetString("LOG_LEVEL", "info")),
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.String(slog.TimeKey, a.Value.Time().Format(time.RFC3339))
			case slog.LevelKey:
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := levelNames[level]
				if !exists {
					levelLabel = level.String()
				}
				return slog.String(slog.LevelKey, levelLabel)
			default:
				return a
			}
		},
	}

	slogHandler := handler.WithRequestID(handler.NewJSONHandler(os.Stdout, opts))

	return &Logger{slog.New(slogHandler)}
}

func (l *Logger) Fatal(msg string) {
	r := slog.NewRecord(time.Now(), Fatal, msg, 0)
	// Get the caller's PC (program counter) and file/line info
	var pcs [1]uintptr
	// Skip 2 frames: runtime.Callers and this function
	runtime.Callers(2, pcs[:])
	r.PC = pcs[0]
	_ = l.Handler().Handle(context.Background(), r)
	os.Exit(1)
}

func (l *Logger) FatalContext(ctx context.Context, msg string) {
	r := slog.NewRecord(time.Now(), Fatal, msg, 0)
	// Get the caller's PC (program counter) and file/line info
	var pcs [1]uintptr
	// Skip 2 frames: runtime.Callers and this function
	runtime.Callers(2, pcs[:])
	r.PC = pcs[0]
	_ = l.Handler().Handle(ctx, r)
	os.Exit(1)
}

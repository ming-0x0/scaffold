package bunslog

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	sloglogger "github.com/ming-0x0/scaffold/internal/infra/logger/slog"
	"github.com/uptrace/bun"
)

type Logger struct {
	Logger            *sloglogger.Logger
	SlowThreshold     time.Duration
	IgnoreNoRowsError bool
	LogLevel          slog.Level
}

func New(opts Logger) *Logger {
	if opts.Logger == nil {
		opts.Logger = sloglogger.New()
	}

	if opts.LogLevel < 0 {
		opts.LogLevel = slog.LevelInfo
	}

	return &opts
}

var _ bun.QueryHook = (*Logger)(nil)

func (l *Logger) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (l *Logger) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	elapsed := time.Since(event.StartTime)

	attrs := []slog.Attr{
		slog.String("durations", elapsed.String()),
		slog.String("sql", event.Query),
	}

	if event.Result != nil {
		rows, err := event.Result.RowsAffected()
		if err != nil {
			attrs = append(attrs, slog.String("rows", "-"))
		} else {
			attrs = append(attrs, slog.Int64("rows", rows))
		}
	}

	if event.Err != nil {
		attrs = append(attrs, slog.String("error", event.Err.Error()))
	}

	// Create a new logger with all attributes
	logger := l.Logger.WithGroup("bun")
	for _, attr := range attrs {
		logger = logger.With(attr.Key, attr.Value.Any())
	}

	switch {
	case event.Err != nil && (!errors.Is(event.Err, sql.ErrNoRows) || !l.IgnoreNoRowsError):
		if l.LogLevel >= sloglogger.Error {
			logger.ErrorContext(ctx, "SQL Query failed")
		}
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold:
		if l.LogLevel >= sloglogger.Warn {
			logger.WarnContext(ctx, "Performed SLOW SQL Query")
		}
	default:
		if l.LogLevel >= sloglogger.Info {
			logger.InfoContext(ctx, "Performed SQL Query")
		}
	}
}

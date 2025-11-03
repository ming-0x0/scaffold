package bunslog

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/uptrace/bun"
)

type Logger struct {
	Logger            *slog.Logger
	SlowThreshold     time.Duration
	IgnoreNoRowsError bool
	LogLevel          slog.Level
}

func New(opts Logger) *Logger {
	if opts.Logger == nil {
		opts.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			AddSource:   true,
			ReplaceAttr: nil,
		}))
	}

	if opts.LogLevel == 0 {
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

	rows, err := event.Result.RowsAffected()
	if err != nil {
		attrs = append(attrs, slog.String("rows", "-"))
	} else {
		attrs = append(attrs, slog.Int64("rows", rows))
	}

	if event.Err != nil {
		attrs = append(attrs, slog.String("error", event.Err.Error()))
	}

	// Create a new logger with all attributes
	logger := l.Logger
	for _, attr := range attrs {
		logger = logger.With(attr.Key, attr.Value.Any())
	}

	switch {
	case event.Err != nil && (!errors.Is(event.Err, sql.ErrNoRows) || !l.IgnoreNoRowsError):
		if l.LogLevel >= sloglogger.Error {
			logger.ErrorContext(ctx, "SQL Query failed", "query", event.Query, "duration", elapsed, "rows", rows, "error", event.Err)
		}
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold:
		if l.LogLevel >= sloglogger.Warn {
			logger.WarnContext(ctx, "Performed SLOW SQL Query", "query", event.Query, "duration", elapsed, "rows", rows)
		}
	default:
		if l.LogLevel >= sloglogger.Info {
			logger.InfoContext(ctx, "Performed SQL Query", "query", event.Query, "duration", elapsed, "rows", rows)
		}
	}
}

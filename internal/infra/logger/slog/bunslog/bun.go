package bunslog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	sloglogger "github.com/ming-0x0/scaffold/internal/infra/logger/slog"
	"github.com/sirupsen/logrus"
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

	if opts.LogLevel > 0 {
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

	fields := logrus.Fields{}
	fields["durations"] = elapsed.String()

	fields["sql"] = event.Query

	if len(event.QueryArgs) > 0 {
		args := make([]string, 0, len(event.QueryArgs))
		for _, arg := range event.QueryArgs {
			args = append(args, fmt.Sprintf("%v", arg))
		}
		fields["args"] = args
	}

	if event.Err != nil {
		fields["error"] = event.Err.Error()
	}

	switch {
	case event.Err != nil && (!errors.Is(event.Err, sql.ErrNoRows) || !l.IgnoreNoRowsError):
		if l.LogLevel >= sloglogger.Error {
			l.Logger.WithContext(ctx).WithFields(fields).Error("SQL Query failed")
		}
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold:
		if l.LogLevel >= slog.WarnLevel {
			l.Logger.WithContext(ctx).WithFields(fields).Warn("Performed SLOW SQL Query")
		}
	default:
		if l.LogLevel >= logrus.InfoLevel {
			l.Logger.WithContext(ctx).WithFields(fields).Info("Performed SQL Query")
		}
	}
}

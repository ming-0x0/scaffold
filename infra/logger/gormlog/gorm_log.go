package gormlog

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/ming-0x0/scaffold/infra/logger"
	"github.com/ming-0x0/scaffold/infra/logger/hook"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Logger based on logrus, but compatible with gorm
type Logger struct {
	Logger                    *logrus.Entry
	LogLevel                  gormLog.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	FileWithLineNumField      string
}

// New creates a new GormLogger instance with the provided options.
// It initializes the logger if not provided and sets the log level to Silent if not specified.
func New(
	opts Logger,
) *Logger {
	if opts.Logger == nil {
		opts.Logger = logrus.NewEntry(logger.New())
	}

	if opts.LogLevel == 0 {
		opts.LogLevel = gormLog.Silent
	}
	return &opts
}

// LogMode sets the log level for the GormLogger and returns the updated logger instance.
func (l *Logger) LogMode(level gormLog.LogLevel) gormLog.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info logs an informational message
func (l *Logger) Info(ctx context.Context, s string, args ...any) {
	if l.LogLevel >= gormLog.Info {
		l.Logger.WithContext(ctx).Infof(s, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, s string, args ...any) {
	if l.LogLevel >= gormLog.Warn {
		l.Logger.WithContext(ctx).Warnf(s, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, s string, args ...any) {
	if l.LogLevel >= gormLog.Error {
		l.Logger.WithContext(ctx).Errorf(s, args...)
	}
}

// Trace logs detailed information about a database query execution, including SQL, duration, and potential errors.
// We want the SQL logs with the info level, while it's defined as trace by gorm
// It checks the log level to determine whether to log the information and formats the log message accordingly.
// Additionally, it handles cases such as slow queries and error messages related to database operations.
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLog.Silent {
		return
	}

	fields := logrus.Fields{}
	fields[l.FileWithLineNumField] = filepath.Base(utils.FileWithLineNum())

	sql, rows := fc()
	if rows == -1 {
		fields["rows"] = "-"
	} else {
		fields["rows"] = rows
	}
	fields["sql"] = sql

	elapsed := time.Since(begin)
	fields["durations"] = elapsed.String()

	if requestID, ok := ctx.Value(hook.RequestIDField).(string); ok {
		fields[hook.RequestIDField] = requestID
	}

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) && l.LogLevel >= gormLog.Error:
		fields["error"] = err.Error()
		l.Logger.WithContext(ctx).WithFields(fields).Error("SQL Query failed")
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormLog.Warn:
		l.Logger.WithContext(ctx).WithFields(fields).Warn("Performed SLOW SQL Query")
	case l.LogLevel == gormLog.Info:
		l.Logger.WithContext(ctx).WithFields(fields).Info("Performed SQL Query")
	}
}

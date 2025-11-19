package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	sloglogger "github.com/ming-0x0/scaffold/internal/infra/logger/slog"
	bunslog "github.com/ming-0x0/scaffold/internal/infra/logger/slog/bun"
	"github.com/ming-0x0/scaffold/pkg/env"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewPostgreSQLDB(
	ctx context.Context,
	logger *sloglogger.Logger,
) (*bun.DB, error) {
	// log that we are opening the connection pool
	logger.InfoContext(ctx, "Opening the PostgreSQL DB connection pool")

	config := &pgdriver.Config{
		Network:  "tcp",
		Addr:     fmt.Sprintf("%s:%d", env.GetString("DB_HOST", "localhost"), env.GetInt("DB_PORT", 5432)),
		User:     env.GetString("DB_USER", "postgres"),
		Password: env.GetString("DB_PASS", "postgres"),
		Database: env.GetString("DB_NAME", "postgres"),
	}

	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithConfig(config))), pgdialect.New())
	db.SetMaxOpenConns(env.GetInt("DB_MAX_OPEN_CONNS", 10))
	db.SetMaxIdleConns(env.GetInt("DB_MAX_IDLE_CONNS", 10))
	db.SetConnMaxLifetime(time.Duration(env.GetInt("DB_CONN_MAX_LIFETIME", 1800)) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(env.GetInt("DB_CONN_MAX_IDLE_TIME", 1800)) * time.Second)

	db.AddQueryHook(bunslog.New(bunslog.Logger{
		Logger:            logger,
		SlowThreshold:     200 * time.Millisecond,
		IgnoreNoRowsError: false,
		LogLevel:          sloglogger.Info,
	}))

	// ping the database
	if err := db.Ping(); err != nil {
		logger.ErrorContext(ctx, "Error while pinging the PostgreSQL DB connection pool", slog.Any("error", err))
		return nil, err
	}

	return db, nil
}

// ClosePostgreSQLDB closes the database connection pool
func ClosePostgreSQLDB(
	ctx context.Context,
	db *bun.DB, // database instance
	logger *sloglogger.Logger, // logger instance
) {
	// log that we are closing the connection pool
	logger.InfoContext(ctx, "Closing the PostgreSQL DB connection pool")

	// close the connection pool
	if err := db.Close(); err != nil {
		logger.ErrorContext(ctx, "Error while closing the PostgreSQL DB connection pool", slog.Any("error", err))
		return
	}
}

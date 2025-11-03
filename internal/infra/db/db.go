package db

import (
	sloglogger "github.com/ming-0x0/scaffold/internal/infra/logger/slog"
	"github.com/uptrace/bun"
)

func New(
	db *bun.DB,
	logger *sloglogger.Logger,
) error {
	logger.Info("Initializing database connection")
	return nil
}

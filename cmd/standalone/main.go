package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/ming-0x0/scaffold/internal/infra/db"
	sloglogger "github.com/ming-0x0/scaffold/internal/infra/logger/slog"
)

func main() {
	logger := sloglogger.New()

	logger.Fatal("hello world")

	_ = db.New(nil, logger)
}

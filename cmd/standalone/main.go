package main

import (
	"context"

	sloglogger "github.com/ming-0x0/scaffold/internal/infra/logger/slog"
)

func main() {
	logger := sloglogger.New()

	logger.InfoContext(context.Background(), "hello world")
}

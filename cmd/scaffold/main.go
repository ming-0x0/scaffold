package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/ming-0x0/scaffold/adapter/gateway"
	"github.com/ming-0x0/scaffold/adapter/grpc"
	"github.com/ming-0x0/scaffold/adapter/repository"
	"github.com/ming-0x0/scaffold/infra/db"
	"github.com/ming-0x0/scaffold/infra/logger"
	"github.com/ming-0x0/scaffold/infra/logger/hook"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger
	logger := logger.New()
	logger.AddHook(&hook.RequestIDHook{})

	// Initialize database
	logger.Info("Initializing database...")
	gormDB, err := db.NewDB(logger)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.CloseDB(gormDB, logger)

	repositoryContainer := repository.NewRepositoryContainer(gormDB, logger)

	grpcServer := grpc.New(
		grpc.Config{
			Port:            os.Getenv("GRPC_PORT"),
			ShutdownTimeout: 5 * time.Second,
		},
		repositoryContainer,
		logger,
	)

	gatewayServer := gateway.New(
		gateway.Config{
			GRPCServer:      os.Getenv("GRPC_SERVER"),
			Port:            os.Getenv("GATEWAY_PORT"),
			ShutdownTimeout: 5 * time.Second,
		},
		logger,
	)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// Start both servers concurrently
	wg.Add(2)

	go func() {
		defer wg.Done()
		logger.Info("Starting gRPC server...")
		if err := grpcServer.Run(ctx); err != nil {
			errChan <- err
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		logger.Info("Starting gateway server...")
		if err := gatewayServer.Run(ctx); err != nil {
			errChan <- err
			cancel()
		}
	}()

	// Handle shutdown signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either signal or error
	select {
	case sig := <-signals:
		logger.WithField("signal", sig).Info("Received shutdown signal")
		cancel()
	case err := <-errChan:
		logger.WithError(err).Error("Service error occurred")
		cancel()
	}

	// Wait for both servers to shutdown
	wg.Wait()
	logger.Info("All services stopped gracefully")
}

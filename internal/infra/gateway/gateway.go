package gateway

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	portalv1 "github.com/ming-0x0/scaffold/api/gen/go/portal/v1"
	"github.com/ming-0x0/scaffold/internal/adapter/gateway/annotator"
	"github.com/ming-0x0/scaffold/internal/adapter/gateway/marshaler"
	"github.com/ming-0x0/scaffold/internal/adapter/gateway/middleware"
	"github.com/ming-0x0/scaffold/internal/adapter/gateway/responder"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	GRPCServer      string
	Port            string
	ShutdownTimeout time.Duration
}

type Server struct {
	config    Config
	annotator annotator.AnnotatorInterface
	responder responder.ResponderInterface
	marshaler marshaler.MarshalerInterface
	logger    *logrus.Logger
}

func New(cfg Config, logger *logrus.Logger) *Server {
	return &Server{
		config:    cfg,
		annotator: annotator.New(),
		responder: responder.New(),
		marshaler: marshaler.New(),
		logger:    logger,
	}
}

func (s *Server) createGRPCConn() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		s.config.GRPCServer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.logger.WithError(err).Error("GRPC connection failed")
		return nil, err
	}

	return conn, nil
}

func (s *Server) createMux() *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithErrorHandler(s.responder.RespondError),
		runtime.WithMetadata(s.annotator.AnnotateMetadata),
		runtime.WithForwardResponseOption(s.responder.Respond),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, s.marshaler.NewNilMarshaler()),
	)
}

func (s *Server) createHTTPServer(mux *runtime.ServeMux) *http.Server {
	mainMux := http.NewServeMux()
	mainMux.Handle("/api/", http.StripPrefix("/api", mux))
	handler := middleware.WithRequestID(mainMux)

	return &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func (s *Server) registerServices(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	for _, fn := range []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error{
		portalv1.RegisterPortalAuthHandler,
	} {
		if err := fn(ctx, mux, conn); err != nil {
			return err
		}
	}
	s.logger.Info("All services registered successfully")
	return nil
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.WithFields(logrus.Fields{
		"grpc_server": s.config.GRPCServer,
		"port":        s.config.Port,
	}).Info("Starting gRPC Gateway server")

	conn, err := s.createGRPCConn()
	if err != nil {
		s.logger.WithError(err).Error("GRPC connection failed")
		return err
	}

	mux := s.createMux()

	if err := s.registerServices(ctx, mux, conn); err != nil {
		s.logger.WithError(err).Error("Service registration failed")
		return err
	}

	server := s.createHTTPServer(mux)

	serveErr := make(chan error, 1)
	go func() {
		s.logger.Infof("HTTP server listening on %s", s.config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErr <- err
		}
	}()

	select {
	case err := <-serveErr:
		s.logger.WithError(err).Error("Server runtime error")
		return err
	case <-ctx.Done():
		s.logger.Info("Shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.WithError(err).Error("HTTP server shutdown failed")
		} else {
			s.logger.Info("HTTP server stopped gracefully")
		}

		if err := conn.Close(); err != nil {
			s.logger.WithError(err).Error("failed to close gRPC connection")
		} else {
			s.logger.Info("gRPC connection closed")
		}
		return nil
	}
}

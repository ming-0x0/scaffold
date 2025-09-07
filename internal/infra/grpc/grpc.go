package grpc

import (
	"context"
	"net"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	portalv1 "github.com/ming-0x0/scaffold/api/gen/go/portal/v1"
	"github.com/ming-0x0/scaffold/internal/adapter"
	v1 "github.com/ming-0x0/scaffold/internal/adapter/grpc/handler/portal/v1"
	interceptorAdapter "github.com/ming-0x0/scaffold/internal/adapter/grpc/interceptor"
	"github.com/ming-0x0/scaffold/internal/adapter/grpc/responder"
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type Config struct {
	Port            string
	ShutdownTimeout time.Duration
}

type Server struct {
	config         Config
	server         *grpc.Server
	adapter        domain.AdapterInterface
	interceptor    interceptorAdapter.InterceptorInterface
	errorResponder responder.ErrorResponderInterface
	logger         *logrus.Logger
}

func New(
	cfg Config,
	db *gorm.DB,
	logger *logrus.Logger,
) *Server {
	s := &Server{
		config:         cfg,
		adapter:        adapter.New(db, logger),
		interceptor:    interceptorAdapter.New(),
		errorResponder: responder.NewErrorResponder(),
		logger:         logger,
	}

	serverOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			s.interceptor.InterceptPanic,
			s.interceptor.InterceptContext,
		),
	}

	s.server = grpc.NewServer(serverOpts...)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.server, healthServer)

	s.registerHandlers()

	reflection.Register(s.server)

	return s
}

func (s *Server) registerHandlers() {
	v1Handler := v1.New(s.adapter, s.errorResponder, s.logger)

	registers := []func(){
		func() { portalv1.RegisterPortalAuthServer(s.server, v1.NewAuthHandler(v1Handler)) },
		func() { portalv1.RegisterPortalBannerServer(s.server, v1.NewBannerHandler(v1Handler)) },
	}

	for _, register := range registers {
		register()
	}
}

func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", ":"+s.config.Port)
	if err != nil {
		s.logger.WithError(err).Error("failed to create listener")
		return err
	}

	serveErr := make(chan error, 1)
	go func() {
		if err := s.server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			serveErr <- err
		}
	}()

	select {
	case err := <-serveErr:
		s.logger.WithError(err).Error("Server runtime error")
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		stopped := make(chan struct{})
		go func() {
			s.server.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			s.logger.Info("Server stopped gracefully")
		case <-shutdownCtx.Done():
			s.logger.Warn("Forcing server shutdown after timeout")
			s.server.Stop()
		}
		return nil
	}
}

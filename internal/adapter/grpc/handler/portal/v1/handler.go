package v1

import (
	"github.com/ming-0x0/scaffold/internal/adapter/grpc/responder"
	"github.com/ming-0x0/scaffold/internal/domain"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	adapter        domain.AdapterInterface
	errorResponder responder.ErrorResponderInterface
	logger         *logrus.Logger
}

func New(
	adapter domain.AdapterInterface,
	errorResponder responder.ErrorResponderInterface,
	logger *logrus.Logger,
) *Handler {
	return &Handler{
		adapter:        adapter,
		errorResponder: errorResponder,
		logger:         logger,
	}
}

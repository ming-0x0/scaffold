package responder

import (
	"errors"

	"github.com/ming-0x0/scaffold/pkg/domainerror"
)

type ErrorResponderInterface interface {
	Respond(err error) error
	RespondMsg(msg string, err error) error
	RespondCode(errCode domainerror.ErrorCode, msg string, err error) error
}

type ErrorResponder struct{}

func NewErrorResponder() *ErrorResponder {
	return &ErrorResponder{}
}

func (e *ErrorResponder) Respond(err error) error {
	var domainErr *domainerror.DomainError
	if errors.As(err, &domainErr) {
		return domainErr.GRPCStatus()
	}

	return domainerror.Wrap(domainerror.Internal, err).GRPCStatus()
}

func (e *ErrorResponder) RespondMsg(msg string, err error) error {
	var domainErr *domainerror.DomainError
	if errors.As(err, &domainErr) {
		return domainerror.WrapMsg(domainErr.ErrorCode(), msg, err).GRPCStatus()
	}

	return domainerror.WrapMsg(domainerror.Internal, msg, err).GRPCStatus()
}

func (e *ErrorResponder) RespondCode(errCode domainerror.ErrorCode, msg string, err error) error {
	if err == nil {
		return domainerror.WrapMsg(errCode, msg, errors.New(msg)).GRPCStatus()
	}
	return domainerror.WrapMsg(errCode, msg, err).GRPCStatus()
}

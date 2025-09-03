package responder

import (
	"errors"

	"github.com/ming-0x0/scaffold/shared/domainerror"
)

type ErrorResponderInterface interface {
	Wrap(err error) error
	WrapMsg(msg string, err error) error
	WrapCode(msg string, code domainerror.ErrorCode, errors ...error) error
}

type ErrorResponder struct{}

func NewErrorResponder() *ErrorResponder {
	return &ErrorResponder{}
}

func (e *ErrorResponder) Wrap(err error) error {
	var domainErr *domainerror.DomainError
	if errors.As(err, &domainErr) {
		return domainErr.GRPCStatus()
	}

	return domainerror.Wrap(domainerror.Internal, err).GRPCStatus()
}

func (e *ErrorResponder) WrapMsg(msg string, err error) error {
	var domainErr *domainerror.DomainError
	if errors.As(err, &domainErr) {
		return domainerror.WrapMsg(domainErr.ErrorCode(), msg, err).GRPCStatus()
	}

	return domainerror.WrapMsg(domainerror.Internal, msg, err).GRPCStatus()
}

func (e *ErrorResponder) WrapCode(msg string, code domainerror.ErrorCode, errors ...error) error {
	return domainerror.WrapMsg(code, msg, errors[0]).GRPCStatus()
}

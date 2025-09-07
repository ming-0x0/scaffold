package domainerror

import (
	"github.com/ming-0x0/scaffold/api/gen/go/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate stringer -type=ErrorCode
type ErrorCode int32

const (
	OK ErrorCode = iota
	Canceled
	Unknown
	InvalidArgument
	DeadlineExceeded
	NotFound
	AlreadyExists
	PermissionDenied
	ResourceExhausted
	FailedPrecondition
	Aborted
	OutOfRange
	Unimplemented
	Internal
	Unavailable
	DataLoss
	Unauthenticated
)

type DomainError struct {
	errCode ErrorCode
	msg     string
	err     error
}

func Wrap(
	errCode ErrorCode,
	err error,
) *DomainError {
	return &DomainError{
		errCode: errCode,
		err:     err,
	}
}

func WrapMsg(
	errCode ErrorCode,
	msg string,
	err error,
) *DomainError {
	return &DomainError{
		errCode: errCode,
		msg:     msg,
		err:     err,
	}
}

func (e *DomainError) Error() string {
	if e == nil {
		return ""
	}

	return e.err.Error()
}

func (e *DomainError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.err
}

func (e *DomainError) ErrorCode() ErrorCode {
	if e == nil {
		return OK
	}

	return e.errCode
}

func (e *DomainError) Message() string {
	if e == nil {
		return ""
	}

	if e.msg == "" {
		return e.errCode.String()
	}

	return e.msg
}

func (e *DomainError) GRPCStatus() error {
	st := status.New(codes.Code(e.ErrorCode()), e.Message())

	details := &common.ErrorDetails{
		Details: e.Error(),
	}

	stWithDetails, err := st.WithDetails(details)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}

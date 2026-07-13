package microerror

import "errors"

type Error struct {
	uniqueID string
	code     int
	message  string
	cause    error
	silent   bool
}

func New(uniqueID string, code int, message string) Error {
	return Error{
		uniqueID: uniqueID,
		code:     code,
		message:  message,
	}
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Unwrap() error {
	return e.cause
}

func (e Error) WithCause(err error) Error {
	if err == nil {
		return e
	}

	e.cause = err
	return e
}

func (e Error) LogOnly() Error {
	e.silent = true
	return e
}

func (e Error) IsSilent() bool {
	return e.silent
}

func (e Error) UniqueID() string {
	return e.uniqueID
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func As(err error) (Error, bool) {
	var appErr Error
	if errors.As(err, &appErr) {
		return appErr, true
	}

	return Error{}, false
}

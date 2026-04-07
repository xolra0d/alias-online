package main

type ServiceError string

const (
	ErrInternal          ServiceError = "internal error"
	ErrWrongCredentials  ServiceError = "wrong credentials"
	ErrUserAlreadyExists ServiceError = "user already exists"
	ErrUserNotFound      ServiceError = "user not found"
)

type Error struct {
	SvcError ServiceError // Outer error returned to user
	AppError error        // Inner error returned by inner services (db)
}

func NewError(svcError ServiceError, appError error) *Error {
	return &Error{SvcError: svcError, AppError: appError}
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return string(e.SvcError)
}

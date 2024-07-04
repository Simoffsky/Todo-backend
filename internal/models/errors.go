package models

import (
	"errors"
	"net/http"
)

type Error struct {
	Err        error
	StatusCode int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func NewError(err error, statusCode int) Error {
	return Error{
		Err:        err,
		StatusCode: statusCode,
	}
}

var (
	ErrUserExists      = NewError(errors.New("user already exists"), http.StatusConflict)
	ErrUserNotFound    = NewError(errors.New("user not found"), http.StatusNotFound)
	ErrInvalidPassword = NewError(errors.New("invalid password"), http.StatusUnauthorized)

	ErrTaskNotFound = NewError(errors.New("task not found"), http.StatusNotFound)

	ErrTaskListNotFound = NewError(errors.New("task list not found"), http.StatusNotFound)

	ErrUnauthorized = NewError(errors.New("unauthorized"), http.StatusUnauthorized)
	ErrInvalidToken = NewError(errors.New("invalid token"), http.StatusUnauthorized)
	ErrTokenExpired = NewError(errors.New("token expired"), http.StatusUnauthorized)

	ErrAccessDenied = NewError(errors.New("access denied"), http.StatusForbidden)

	ErrForbidden = NewError(errors.New("forbidden"), http.StatusForbidden)
)

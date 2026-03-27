package utils

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func NewBadRequest(message string, err error) *AppError {
	return &AppError{Code: "bad_request", Message: message, HTTPStatus: http.StatusBadRequest, Err: err}
}

func NewUnauthorized(message string, err error) *AppError {
	return &AppError{Code: "unauthorized", Message: message, HTTPStatus: http.StatusUnauthorized, Err: err}
}

func NewForbidden(message string, err error) *AppError {
	return &AppError{Code: "forbidden", Message: message, HTTPStatus: http.StatusForbidden, Err: err}
}

func NewNotFound(message string, err error) *AppError {
	return &AppError{Code: "not_found", Message: message, HTTPStatus: http.StatusNotFound, Err: err}
}

func NewConflict(message string, err error) *AppError {
	return &AppError{Code: "conflict", Message: message, HTTPStatus: http.StatusConflict, Err: err}
}

func NewInternal(message string, err error) *AppError {
	return &AppError{Code: "internal", Message: message, HTTPStatus: http.StatusInternalServerError, Err: err}
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

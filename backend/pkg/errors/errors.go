package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom error type for our application
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// New creates a new AppError
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// BadRequest returns a 400 Bad Request error
func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

// Unauthorized returns a 401 Unauthorized error
func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}

// Forbidden returns a 403 Forbidden error
func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, message)
}

// NotFound returns a 404 Not Found error
func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message)
}

// InternalServerError returns a 500 Internal Server Error
func InternalServerError(message string) *AppError {
	return New(http.StatusInternalServerError, message)
}

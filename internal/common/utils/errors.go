package utils

import "fmt"

// AppError represents a custom application error
type AppError struct {
	StatusCode int
	Message    string
	Err        error
}

// Error allows AppError to satisfy the standard error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError is a constructor for AppError
func NewAppError(statusCode int, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

package lib

import (
	"fmt"
	"net/http"
)

// APIError represents an error that provides an HTTP status code, message and causes.
type APIError interface {
	Error() string

	Wrap(err error) APIError
}

// apiError is a concrete implementation of the APIError interface.
type apiError struct {
	Message string        `json:"message"`
	Status  int           `json:"status"`
	Causes  []interface{} `json:"causes"`
}

// Error implements the error interface.
func (e apiError) Error() string {
	return fmt.Sprintf("message: %s - status: %d - causes: %v",
		e.Message, e.Status, e.Causes)
}

// Wrap wraps an existing error into an APIError.
func (e apiError) Wrap(err error) APIError {
	if err != nil {
		e.Causes = append(e.Causes, err.Error())
	}
	return e
}

// NewBadRequestError creates a new APIError for bad requests.
//
// Example usage:
//
//	err := NewBadRequestError("invalid input").Wrap(innerErr)
func NewBadRequestError(message string) APIError {
	return apiError{
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// NewNotFoundError creates a new APIError for not found errors.
//
// Example usage:
//
//	err := NewNotFoundError("resource not found")
func NewNotFoundError(message string) APIError {
	return apiError{
		Message: message,
		Status:  http.StatusNotFound,
	}
}

// NewUnauthorizedError creates a new APIError for unauthorized requests.
//
// Example usage:
//
//	err := NewUnauthorizedError("unauthorized")
func NewUnauthorizedError(message string) APIError {
	return apiError{
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

// NewUnexpectedError creates a new APIError for unexpected errors.
//
// Example usage:
//
//	err := NewUnexpectedError("something went wrong")
func NewUnexpectedError(message string) APIError {
	return apiError{
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

// NewInternalServerError creates a new APIError for internal server errors.
//
// Example usage:
//
//	err := NewInternalServerError("internal server error", innerErr)
func NewInternalServerError(message string, err error) APIError {
	result := apiError{
		Message: message,
		Status:  http.StatusInternalServerError,
	}
	if err != nil {
		result.Causes = append(result.Causes, err.Error())
	}
	return result
}

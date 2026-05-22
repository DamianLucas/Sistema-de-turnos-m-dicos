package pkg

import "net/http"

// HTTPError representa un error que conoce el status HTTP que debe devolver.
type HTTPError interface {
	error
	StatusCode() int
}

// AppError implementa HTTPError.
type AppError struct {
	status  int
	message string
}

// Error implementa la interfaz error.
func (e *AppError) Error() string { return e.message }

// StatusCode devuelve el status HTTP asociado.
func (e *AppError) StatusCode() int { return e.status }

//Helpers

func NewBadRequestError(message string) error {
	return &AppError{
		status:  http.StatusBadRequest,
		message: message,
	}
}

func NewNotFoundError(message string) error {
	return &AppError{
		status:  http.StatusNotFound,
		message: message,
	}
}

func NewConflictError(message string) error {
	return &AppError{
		status:  http.StatusConflict,
		message: message,
	}
}

func NewUnauthorizedError(message string) error {
	return &AppError{
		status:  http.StatusUnauthorized,
		message: message,
	}
}

func NewForbiddenError(message string) error {
	return &AppError{
		status:  http.StatusForbidden,
		message: message,
	}
}

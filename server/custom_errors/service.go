package custom_errors

import "errors"

type ServiceError struct {
	error   `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *ServiceError) Unwrap() error {
	return e.error
}

func NewServiceError(message string, status int) *ServiceError {
	return &ServiceError{
		error:   errors.New(message),
		Status:  status,
		Message: message,
	}
}

package rest

import "fmt"

type ApiResponse[T any] struct {
	IsSuccess bool      `json:"is_success"`
	Data      T         `json:"data"`
	Error     *ApiError `json:"error"`
}

type ApiError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("code=%d, message=%s", e.StatusCode, e.Message)
}

func NewApiError(code, message string, statusCode int) *ApiError {
	return &ApiError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

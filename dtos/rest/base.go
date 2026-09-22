package rest

import "fmt"

type ApiResponse[T any] struct {
	IsSuccess bool      `json:"isSuccess"`
	Data      T         `json:"data"`
	Error     *ApiError `json:"error"`
}

type ApiError struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("code=%d, messages=%s", e.StatusCode, e.Message)
}

func NewApiError(code, message string, statusCode int) *ApiError {
	return &ApiError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

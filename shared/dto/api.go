package dto

type ApiError struct {
	Name    string            `json:"name"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type ApiResponse[T any] struct {
	Data   T         `json:"data,omitempty"`
	Status string    `json:"status"`
	Error  *ApiError `json:"error,omitempty"`
}

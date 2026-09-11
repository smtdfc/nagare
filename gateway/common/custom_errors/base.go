package custom_errors

import "errors"

type GatewayError struct {
	error
	Details    string
	Code       string
	StatusCode int
}

func (n *GatewayError) Error() string {
	return n.Details
}

func NewGatewayError(code, details string, statusCode int) *GatewayError {
	return &GatewayError{
		error:      errors.New(details),
		Details:    details,
		Code:       code,
		StatusCode: statusCode,
	}
}

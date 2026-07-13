package custom_errors

import "errors"

type StreamError struct {
	error
}

func NewStreamError(msg string) *StreamError {
	return &StreamError{
		error: errors.New(msg),
	}
}

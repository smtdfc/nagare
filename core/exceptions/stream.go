package exceptions

import "errors"

type StreamException struct {
	error
}

func NewStreamException(msg string) *StreamException {
	return &StreamException{
		error: errors.New(msg),
	}
}

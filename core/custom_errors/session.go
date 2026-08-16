package custom_errors

import "errors"

type SessionError struct {
	error
}

func NewSessionError(msg string) *SessionError {
	return &SessionError{
		error: errors.New(msg),
	}
}

package custom_errors

import "errors"

type DatabaseError struct {
	error
}

func NewDatabaseError(msg string) *DatabaseError {
	return &DatabaseError{
		error: errors.New(msg),
	}
}

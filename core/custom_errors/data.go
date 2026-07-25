package custom_errors

import "errors"

type DataError struct {
	error
}

func NewDataError(msg string) *DataError {
	return &DataError{
		error: errors.New(msg),
	}
}

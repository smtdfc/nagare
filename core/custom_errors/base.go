package custom_errors

import "errors"

type NagareCoreError struct {
	error
	Details string
	Code    string
}

func (n *NagareCoreError) Error() string {
	return n.Details
}

func NewNagareCoreError(code, details string) *NagareCoreError {
	return &NagareCoreError{
		error:   errors.New(details),
		Details: details,
		Code:    code,
	}
}

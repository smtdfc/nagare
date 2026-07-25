package custom_errors

import "errors"

type ToolError struct {
	error
	Name string
}

func NewToolError(msg string, name string) *ToolError {
	return &ToolError{
		error: errors.New(msg),
		Name:  name,
	}
}

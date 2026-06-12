package exceptions

import "errors"

type ToolException struct {
	error
	tool string
}

func NewToolException(msg string, toolName string) *ToolException {
	return &ToolException{
		error: errors.New(msg),
		tool:  toolName,
	}
}

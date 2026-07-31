package custom_errors

import "errors"

type AgentError struct {
	error
}

func NewAgentError(msg string) *AgentError {
	return &AgentError{
		error: errors.New(msg),
	}
}

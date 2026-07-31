package custom_errors

import "errors"

type LLMProviderError struct {
	error
}

func NewLLMProviderError(msg string) *LLMProviderError {
	return &LLMProviderError{
		error: errors.New(msg),
	}
}

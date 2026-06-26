package exceptions

import "errors"

type LLMProviderException struct {
	error
	Provider string
}

func NewLLMProviderException(msg string, provider string) *LLMProviderException {
	return &LLMProviderException{
		error:    errors.New(msg),
		Provider: provider,
	}
}

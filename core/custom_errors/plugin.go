package custom_errors

import "errors"

type PluginError struct {
	error
}

func NewPluginError(msg string) *PluginError {
	return &PluginError{
		error: errors.New(msg),
	}
}

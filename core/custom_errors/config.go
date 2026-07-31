package custom_errors

import "errors"

type ConfigError struct {
	error
}

func NewConfigError(msg string) *ConfigError {
	return &ConfigError{
		error: errors.New(msg),
	}
}

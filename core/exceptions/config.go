package exceptions

import "errors"

type ConfigException struct {
	error
}

func NewConfigException(msg string) *ConfigException {
	return &ConfigException{
		error: errors.New(msg),
	}
}

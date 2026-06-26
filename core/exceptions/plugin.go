package exceptions

import "errors"

type PluginException struct {
	error
	Name string
}

func NewPluginException(msg string, name string) *PluginException {
	return &PluginException{
		error: errors.New(msg),
		Name:  name,
	}
}

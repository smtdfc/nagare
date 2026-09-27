package tool

import (
	context2 "github.com/smtdfc/nagare/core/context"
)

type Tool interface {
	GetName() string
	GetDescription() string
	GetArgsSchema() string
	GetBindings() Bindings
	GetCategories() []string
	WithBindings(bindings Bindings) Tool
	Execute(ctx *context2.ExecuteContext, args string) (string, error)
}

type ListTool []Tool

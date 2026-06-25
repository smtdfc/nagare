package tool

import (
	"context"
)

type Tool interface {
	GetName() string
	Execute(context.Context, string) (string, error)
	GetArgumentsSchema() string
	GetDesc() string
	GetType() ToolType
}

type ListTool []Tool

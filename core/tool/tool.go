package tool

import "context"

type Tool interface {
	GetName() string
	GetDescription() string
	GetArgsSchema() string
	Execute(ctx context.Context, args string) (string, error)
}

type ListTool []Tool

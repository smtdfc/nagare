package tool

import (
	context2 "github.com/smtdfc/nagare/core/context"
)

type Bindings interface {
	RefreshTask(ctx *context2.ExecuteContext)
	CreateTask(ctx *context2.ExecuteContext, sessionID string, name string, prompt string, triggerBy string, repeat bool, repeatRule string, startTime string, endTime string) (string, error)
	FindToolsByKeywords(ctx *context2.ExecuteContext, keywords []string) ([]ToolMetadata, error)
	CallTool(ctx *context2.ExecuteContext, toolName string, args string) *Result
}

type Tool interface {
	GetName() string
	GetDescription() string
	GetArgsSchema() string
	GetBindings() Bindings
	WithBindings(bindings Bindings) Tool
	Execute(ctx *context2.ExecuteContext, args string) (string, error)
}

type ListTool []Tool

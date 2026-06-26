package domains

import "context"

type MiddlewareArg interface {
	UseToolsForNextTurn(ListTool)
}

type MiddlewareFunc func(arg MiddlewareArg)

type AgentContext interface {
	context.Context
	AfterToolCall(MiddlewareFunc, bool)
	GetToolByCategories([]string) ListTool
}

package domains

import "context"

type Context interface {
	context.Context
	ExecuteToolCalls(toolCall *ToolCall) *ToolResult
}

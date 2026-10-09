package tool

import context2 "github.com/smtdfc/nagare/core/context"

type Bindings interface {
	RefreshTask(ctx *context2.ExecuteContext)
	CreateTask(ctx *context2.ExecuteContext, sessionID string, name string, prompt string, triggerBy string, repeat bool, repeatRule string, startTime string, endTime string) (string, error)
	FindToolsByCategories(ctx *context2.ExecuteContext, categories []string) ([]Metadata, error)
	CallTool(ctx *context2.ExecuteContext, toolName string, args string) *Result
}

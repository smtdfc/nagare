package context

import (
	"context"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool"
)

type ExecuteContext struct {
	context.Context
	ToolMgr *tool.ToolManager
}

func (c *ExecuteContext) ExecuteToolCalls(toolCall *domains.ToolCall) *domains.ToolResult {
	result := domains.NewToolResult(toolCall.CallID)
	resultJson, err := c.ToolMgr.CallTool(c, toolCall.Name, toolCall.Args)
	if err != nil {
		return result.Failed(err.Error())
	}

	return result.Success(resultJson)
}

func NewExecuteContext(toolMgr *tool.ToolManager) *ExecuteContext {
	return &ExecuteContext{
		Context: context.Background(),
		ToolMgr: toolMgr,
	}
}

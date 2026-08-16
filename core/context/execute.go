package context

import (
	"context"
	"time"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool"
)

type ExecuteContext struct {
	context.Context
	ToolMgr *tool.ToolManager
}

func (c *ExecuteContext) ExecuteToolCalls(toolCall *domains.ToolCall) *domains.ToolResult {
	startTime := time.Now()
	logger.Logger.Info("Executing tool call", "name", toolCall.Name, "args", toolCall.Args)
	result := domains.NewToolResult(toolCall.CallID)
	resultJson, err := c.ToolMgr.CallTool(c, toolCall.Name, toolCall.Args)
	if err != nil {
		logger.Logger.Error("Failed to execute tool call", "name", toolCall.Name, "args", toolCall.Args, "error", err)
		return result.Failed(err.Error())
	}

	endTime := time.Now()
	logger.Logger.Info("Tool call executed successfully", "name", toolCall.Name, "args", toolCall.Args, "duration", endTime.Sub(startTime))
	return result.Success(resultJson)
}

func NewExecuteContext(toolMgr *tool.ToolManager) *ExecuteContext {
	return &ExecuteContext{
		Context: context.Background(),
		ToolMgr: toolMgr,
	}
}

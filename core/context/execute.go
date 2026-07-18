package context

import (
	"context"

	"github.com/smtdfc/nagare/core/tool"
)

type ExecuteContext struct {
	context.Context
	ToolMgr *tool.ToolManager
}

func NewExecuteContext(toolMgr *tool.ToolManager) *ExecuteContext {
	return &ExecuteContext{
		Context: context.Background(),
		ToolMgr: toolMgr,
	}
}

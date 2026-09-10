package manager

import (
	"context"
	"errors"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool"
)

type ToolManager struct {
	toolMap map[string]tool.Tool
	logger  *logger.BaseLogger
}

func (t *ToolManager) GetListTool() tool.ListTool {
	list := make(tool.ListTool, 0, len(t.toolMap))
	for _, item := range t.toolMap {
		list = append(list, item)
	}

	return list
}

func (t *ToolManager) Call(ctx context.Context, toolCall *tool.ToolCall) *tool.Result {
	toolResultBuilder := tool.NewToolResultBuilder(toolCall.CallID, toolCall.Name)
	calledTool, isExist := t.toolMap[toolCall.Name]
	if !isExist {
		return toolResultBuilder.Failure(errors.New("tool doesn't exist")).Build()
	}

	result, err := calledTool.Execute(ctx, toolCall.Args)
	if err != nil {
		return toolResultBuilder.Failure(err).Build()
	}

	return toolResultBuilder.Success(result).Build()
}

// @Injectable
func NewToolManager(logger *logger.BaseLogger) *ToolManager {
	return &ToolManager{
		toolMap: map[string]tool.Tool{},
		logger:  logger.With("module", "tool-manager"),
	}
}

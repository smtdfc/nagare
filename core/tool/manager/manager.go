package manager

import (
	"context"
	"errors"
	"sync"

	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/tool/registry"
)

type ToolManager struct {
	mu         sync.RWMutex
	cachedList tool.ListTool
	logger     *logger.BaseLogger
}

func (t *ToolManager) GetListTool() tool.ListTool {
	t.mu.RLock()
	if t.cachedList != nil {
		t.mu.RUnlock()
		return t.cachedList
	}
	t.mu.RUnlock()

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.cachedList != nil {
		return t.cachedList
	}

	list := make(tool.ListTool, 0, len(registry.Registry))
	for _, item := range registry.Registry {
		list = append(list, item)
	}

	t.cachedList = list
	return t.cachedList
}

func (t *ToolManager) Call(ctx context.Context, toolCall *tool.ToolCall) *tool.Result {
	toolResultBuilder := tool.NewToolResultBuilder(toolCall.CallID, toolCall.Name)
	calledTool, isExist := registry.Registry[toolCall.Name]
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
		logger: logger.With("module", "tool-manager"),
	}
}

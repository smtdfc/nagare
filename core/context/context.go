package context

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/smtdfc/nagare/core/exceptions"
	"github.com/smtdfc/nagare/core/tool"
	nagare_logger "github.com/smtdfc/nagare/shared/logger"
)

type ExecuteContext struct {
	context.Context
	Tools  tool.ToolMap
	logger *slog.Logger
}

func (e *ExecuteContext) CallTool(name string, args string) (string, error) {
	start := time.Now()
	tool, ok := e.Tools[name]
	if !ok {
		e.logger.Error(fmt.Sprintf("Tool %s not found ", name))
		return "", exceptions.NewToolException(fmt.Sprintf("tool not found: %s", name), name)
	}

	r, err := tool.Execute(e, args)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Tool return with error: %s", err))
		return "", err
	}
	elapsed := time.Since(start)
	e.logger.Info(
		fmt.Sprintf("Tool %s executed", name),
		"duration", elapsed.String(),
		"tool_name", name,
	)
	return r, nil
}

func NewExecuteContext(ctx context.Context, listTool tool.ListTool) ExecuteContext {
	toolMap := tool.ToolMap{}
	for _, tool := range listTool {
		toolMap[tool.GetName()] = tool
	}

	return ExecuteContext{
		Context: ctx,
		Tools:   toolMap,
		logger:  nagare_logger.GetLogger("Execute context"),
	}
}

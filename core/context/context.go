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
	logger *slog.Logger
}

func (e *ExecuteContext) CallTool(name string, args string) (string, error) {
	start := time.Now()
	tool, ok := tool.GlobalToolRegistry.GetByName(name)
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

func NewExecuteContext(ctx context.Context) ExecuteContext {
	return ExecuteContext{
		Context: ctx,
		logger:  nagare_logger.GetLogger("Execute context"),
	}
}

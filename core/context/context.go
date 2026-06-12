package context

import (
	"context"
	"fmt"

	"github.com/smtdfc/nagare/core/exceptions"
	"github.com/smtdfc/nagare/core/tool"
)

type ExecuteContext struct {
	context.Context
	Tools tool.ToolMap
}

func (e *ExecuteContext) CallTool(name string, args string) (string, error) {
	tool, ok := e.Tools[name]
	if !ok {
		return "", exceptions.NewToolException(fmt.Sprintf("tool not found: %s", name), name)
	}

	r, err := tool.Execute(args)
	if err != nil {
		return "", err
	}
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
	}
}

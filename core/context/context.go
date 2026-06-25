package context

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/exceptions"
	"github.com/smtdfc/nagare/core/tool"
	nagare_logger "github.com/smtdfc/nagare/shared/logger"
)

type MiddlewareKind int

const (
	AFTER_TOOL_CALL MiddlewareKind = iota
)

type Middleware struct {
	Kind     MiddlewareKind
	Func     domains.MiddlewareFunc
	CallOnce bool
}

type ExecuteContext struct {
	context.Context
	Middlewares []Middleware
	logger      *slog.Logger
}

func (e *ExecuteContext) AfterToolCall(fn domains.MiddlewareFunc, once bool) {
	e.Middlewares = append(e.Middlewares, Middleware{
		Kind:     AFTER_TOOL_CALL,
		Func:     fn,
		CallOnce: once,
	})
}

func (e *ExecuteContext) CallMiddlewares(kind MiddlewareKind, arg domains.MiddlewareArg) {
	newMiddlewares := make([]Middleware, 0, len(e.Middlewares))
	for _, mid := range e.Middlewares {
		if mid.Kind == kind {
			mid.Func(arg)

			if mid.CallOnce {
				continue
			}
		}
		newMiddlewares = append(newMiddlewares, mid)
	}
	e.Middlewares = newMiddlewares
}

func (e *ExecuteContext) CallTool(name string, args string) (*domains.ToolCallResult, error) {
	start := time.Now()
	tool, ok := tool.GlobalToolRegistry.GetByName(name)

	if !ok {
		e.logger.Error(fmt.Sprintf("Tool %s not found ", name))
		return nil, exceptions.NewToolException(fmt.Sprintf("tool not found: %s", name), name)
	}

	r, err := tool.Execute(e, args)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Tool return with error: %s", err))
		return nil, err
	}
	elapsed := time.Since(start)
	e.logger.Info(
		fmt.Sprintf("Tool %s executed", name),
		"duration", elapsed.String(),
		"tool_name", name,
	)
	return &domains.ToolCallResult{
		Result: r,
		Tool:   tool,
	}, nil
}

func NewExecuteContext(ctx context.Context) ExecuteContext {
	return ExecuteContext{
		Context:     ctx,
		logger:      nagare_logger.GetLogger("Execute context"),
		Middlewares: make([]Middleware, 0),
	}
}

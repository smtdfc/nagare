package context

import "context"

type ExecuteContext struct {
	context.Context
}

func NewExecuteContext() *ExecuteContext {
	return &ExecuteContext{
		Context: context.Background(),
	}
}

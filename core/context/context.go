package context

import "context"

type ExecuteContext struct {
	context.Context

	SessionID string
}

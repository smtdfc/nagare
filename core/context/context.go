package context

import (
	"context"

	"github.com/smtdfc/nagare/core/message"
)

type ExecuteContext struct {
	context.Context
	MessageChannel message.WriteOnlyChannel
	SessionID      string
}

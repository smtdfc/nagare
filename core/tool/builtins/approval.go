package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

// requireApproval asks the human operator to confirm a tool action before it
// touches the system. It fails closed: when no approval channel is available
// the action is refused instead of being executed silently.
func requireApproval(ctx *context.ExecuteContext, bindings tool.Bindings, title string) error {
	if bindings == nil {
		return fmt.Errorf("approval is required but no approval channel is available")
	}

	if !bindings.RequestToolApproval(ctx, title) {
		return fmt.Errorf("action denied by user: %s", title)
	}

	return nil
}

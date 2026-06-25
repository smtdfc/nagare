package declarations

import (
	"context"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/cross-platform/system"
)

type PowerArgs struct {
	Action string `json:"action" jsonschema:"description=Action to perform: shutdown, reboot, sleep"`
}

var power_control = tool.DeclareTool(
	"power_control",
	"Control system power states: shutdown, reboot, or sleep.",
	func(ctx context.Context, args PowerArgs) (any, error) {
		err := system.ExecutePowerCommand(args.Action)
		if err != nil {
			return map[string]any{
				"action":  args.Action,
				"success": false,
				"error":   err.Error(),
			}, nil
		}

		return map[string]any{
			"action":  args.Action,
			"success": true,
		}, nil
	},
	tool.STATIC_TOOL,
	tool.NO_GROUP,
)

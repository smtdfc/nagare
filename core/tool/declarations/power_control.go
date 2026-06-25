package declarations

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/cross-platform/system"
)

type PowerArgs struct {
	Action string `json:"action" jsonschema:"description=Action to perform: shutdown, reboot, sleep, lock_screen"`
}

var power_control = tool.DeclareTool(
	"power_control",
	"Control system power states: shutdown, reboot,lock_screen, or sleep.",
	func(ctx domains.AgentContext, args PowerArgs) (any, error) {
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
	domains.STATIC_TOOL,
	domains.NO_GROUP,
)

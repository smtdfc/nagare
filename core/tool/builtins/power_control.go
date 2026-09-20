package declarations

import (
	"fmt"
	"strings"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type PowerControlInput struct {
	Action string `json:"action"` // "shutdown", "restart", "suspend", "hibernate", "logout"
}

type PowerControlOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var PowerControlTool = tool.DefineTool(
	"power_control_tool",
	"Cross-platform control for system power actions: shutdown, restart, suspend, hibernate, or logout",
	func(ctx *context.ExecuteContext, args *PowerControlInput, _ tool.Bindings) (*PowerControlOutput, error) {
		if args == nil || args.Action == "" {
			return nil, fmt.Errorf("action is required (shutdown, restart, suspend, hibernate, logout)")
		}

		action := strings.ToLower(strings.TrimSpace(args.Action))
		powerCtrl := system.NewPowerControl()
		var err error

		switch action {
		case "shutdown":
			err = powerCtrl.Shutdown()
		case "restart", "reboot":
			err = powerCtrl.Restart()
		case "suspend", "sleep":
			err = powerCtrl.Suspend()
		case "hibernate":
			err = powerCtrl.Hibernate()
		case "logout":
			err = powerCtrl.Logout()
		default:
			return &PowerControlOutput{
				Success: false,
				Message: fmt.Sprintf("Unknown action: '%s'. Supported: shutdown, restart, suspend, hibernate, logout", args.Action),
			}, nil
		}

		if err != nil {
			return &PowerControlOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to execute action '%s': %v", action, err),
			}, nil
		}

		return &PowerControlOutput{
			Success: true,
			Message: fmt.Sprintf("Successfully executed system action: %s", action),
		}, nil
	},
)

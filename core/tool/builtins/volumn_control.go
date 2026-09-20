package declarations

import (
	"fmt"
	"strings"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type VolumeControlInput struct {
	Action string `json:"action"` // "set", "get", "mute", "unmute", "increment", "decrement"
	Value  int    `json:"value"`  // Giá trị tương ứng cho action (ví dụ: 50 cho set, 10 cho tăng/giảm)
}

type VolumeControlOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var VolumeControlTool = tool.DefineTool(
	"volume_control_tool",
	"Cross-platform control for system volume: set, get, mute, unmute, increment, or decrement",
	func(ctx *context.ExecuteContext, args *VolumeControlInput, _ tool.Bindings) (*VolumeControlOutput, error) {
		if args == nil || args.Action == "" {
			return nil, fmt.Errorf("action is required (set, get, mute, unmute, increment, decrement)")
		}

		action := strings.ToLower(strings.TrimSpace(args.Action))
		volCtrl := system.NewVolumeControl()
		var err error

		switch action {
		case "set":
			err = volCtrl.Set(args.Value)
		case "get":
			var current int
			current, err = volCtrl.Get()
			return &VolumeControlOutput{
				Success: false,
				Message: fmt.Sprintf("Current volume: %d", current),
			}, nil
		case "mute":
			err = volCtrl.Mute()
		case "increment":
			err = volCtrl.Increment(args.Value)
		case "decrement":
			err = volCtrl.Decrement(args.Value)
		default:
			return &VolumeControlOutput{
				Success: false,
				Message: fmt.Sprintf("Unknown action: '%s'. Supported: set, get, mute, unmute, increment, decrement", args.Action),
			}, nil
		}

		if err != nil {
			return &VolumeControlOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to execute action '%s': %v", action, err),
			}, nil
		}

		return &VolumeControlOutput{
			Success: true,
			Message: fmt.Sprintf("Successfully executed volume action: %s", action),
		}, nil
	},
)

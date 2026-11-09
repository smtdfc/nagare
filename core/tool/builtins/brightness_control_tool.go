package declarations

import (
	"fmt"
	"strings"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type BrightnessControlToolInput struct {
	Action string `json:"action" jsonschema_description:"Brightness action to perform. Supported actions include 'set', 'get', 'increment', and 'decrement'."`
	Value  int    `json:"value,omitempty" jsonschema_description:"Brightness level used by the action as a percentage between 0 and 100. For 'set' this is the target level; for 'increment' and 'decrement' it is the change amount."`
}

type BrightnessControlToolOutput struct {
	Success bool   `json:"success" jsonschema_description:"Indicates whether the action completed successfully."`
	Level   int    `json:"level,omitempty" jsonschema_description:"The brightness level after the action, as a percentage between 0 and 100."`
	Message string `json:"message" jsonschema_description:"Human readable result of the action."`
}

var BrightnessControlTool = tool.DefineTool(
	"brightness_control_tool",
	"Cross-platform control for display brightness: get, set, increment, or decrement the brightness of the active display as a percentage between 0 and 100",
	func(ctx *context.ExecuteContext, args *BrightnessControlToolInput, _ tool.Bindings) (*BrightnessControlToolOutput, error) {
		if args == nil || args.Action == "" {
			return nil, fmt.Errorf("action is required (set, get, increment, decrement)")
		}

		action := strings.ToLower(strings.TrimSpace(args.Action))
		control := system.NewBrightnessControl()
		var err error

		switch action {
		case "get":
			level, getErr := control.Get()
			if getErr != nil {
				return &BrightnessControlToolOutput{Success: false, Message: fmt.Sprintf("Failed to read brightness: %v", getErr)}, nil
			}
			return &BrightnessControlToolOutput{Success: true, Level: level, Message: fmt.Sprintf("Current brightness: %d", level)}, nil
		case "set":
			err = control.Set(args.Value)
		case "increment":
			var level int
			level, err = control.Increment(args.Value)
			if err == nil {
				return &BrightnessControlToolOutput{Success: true, Level: level, Message: fmt.Sprintf("Brightness increased to %d", level)}, nil
			}
		case "decrement":
			var level int
			level, err = control.Decrement(args.Value)
			if err == nil {
				return &BrightnessControlToolOutput{Success: true, Level: level, Message: fmt.Sprintf("Brightness decreased to %d", level)}, nil
			}
		default:
			return &BrightnessControlToolOutput{
				Success: false,
				Message: fmt.Sprintf("Unknown action: '%s'. Supported: set, get, increment, decrement", args.Action),
			}, nil
		}

		if err != nil {
			return &BrightnessControlToolOutput{Success: false, Message: fmt.Sprintf("Failed to execute action '%s': %v", action, err)}, nil
		}

		return &BrightnessControlToolOutput{Success: true, Level: args.Value, Message: fmt.Sprintf("Successfully executed brightness action: %s", action)}, nil
	},
	[]string{tool.SystemCategory},
)

package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type WindowControlInput struct {
	InstanceID string `json:"instance_id" jsonschema_description:"Unique tracking identifier of the target application instance obtained from get_app_instance_tool."`
	Action     string `json:"action" jsonschema_description:"Window action to perform. Supported actions: 'activate' (bring to foreground and give keyboard focus), 'minimize' (minimize to dock/taskbar), 'maximize' (maximize window), 'restore' (restore from minimized/maximized to normal size), 'close' (close the window)."`
}

type WindowControlOutput struct {
	Success bool `json:"success" jsonschema_description:"Indicates whether the window operation succeeded."`
}

var WindowControlTool = tool.DefineTool(
	"window_control_tool",
	"Manage desktop application windows: focus and bring to front (activate), minimize, maximize, restore, or close.",
	func(ctx *context.ExecuteContext, args *WindowControlInput, bindings tool.Bindings) (*WindowControlOutput, error) {
		if args == nil || args.InstanceID == "" || args.Action == "" {
			return nil, fmt.Errorf("instance_id and action cannot be empty")
		}

		err := bindings.GetComputerUseOrchestrator().WindowControl(ctx, args.InstanceID, args.Action)
		if err != nil {
			return nil, err
		}

		return &WindowControlOutput{Success: true}, nil
	},
	[]string{tool.ComputerUse},
)

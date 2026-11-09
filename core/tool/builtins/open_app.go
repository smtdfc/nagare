package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type OpenAppInput struct {
	Name string   `json:"name" jsonschema_description:"Name or executable path of the application to launch (e.g. 'gnome-calculator', 'firefox', 'notepad', 'code')."`
	Args []string `json:"args,omitempty" jsonschema_description:"Optional command-line arguments to pass to the launched application."`
}

type OpenAppOutput struct {
	PID     *uint32 `json:"pid,omitempty" jsonschema_description:"Process ID (PID) of the launched application process if available."`
	AppName string  `json:"app_name" jsonschema_description:"Name of the application that was requested to launch."`
	Success bool    `json:"success" jsonschema_description:"Indicates whether the application was successfully spawned."`
}

var OpenAppTool = tool.DefineTool(
	"open_app_tool",
	"Launch a new application process by its executable name or path with optional command-line arguments. Requires user approval.",
	func(ctx *context.ExecuteContext, args *OpenAppInput, bindings tool.Bindings) (*OpenAppOutput, error) {
		if args == nil || args.Name == "" {
			return nil, fmt.Errorf("app name cannot be empty")
		}

		title := fmt.Sprintf("Allow the agent to launch application %q", args.Name)
		if len(args.Args) > 0 {
			title = fmt.Sprintf("Allow the agent to launch application %q with args %v", args.Name, args.Args)
		}
		if err := requireApproval(ctx, bindings, title); err != nil {
			return nil, err
		}

		pid, err := bindings.GetComputerUseOrchestrator().OpenApp(ctx, args.Name, args.Args)
		if err != nil {
			return nil, err
		}

		return &OpenAppOutput{
			PID:     pid,
			AppName: args.Name,
			Success: true,
		}, nil
	},
	[]string{tool.ComputerUse},
)

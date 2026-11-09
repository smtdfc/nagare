package declarations

import (
	"github.com/smtdfc/nagare/core/computer_use"
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type ListRunningAppsInput struct{}

type ListRunningAppsOutput struct {
	Apps  []computer_use.AppInfo `json:"apps" jsonschema_description:"List of currently running GUI applications with their names and process IDs."`
	Total int                    `json:"total" jsonschema_description:"Total number of running applications discovered."`
}

var ListRunningAppsTool = tool.DefineTool(
	"list_running_apps_tool",
	"List all currently running desktop applications accessible through the accessibility system.",
	func(ctx *context.ExecuteContext, args *ListRunningAppsInput, bindings tool.Bindings) (*ListRunningAppsOutput, error) {
		apps, err := bindings.GetComputerUseOrchestrator().ListRunningApps(ctx)
		if err != nil {
			return nil, err
		}

		return &ListRunningAppsOutput{
			Apps:  apps,
			Total: len(apps),
		}, nil
	},
	[]string{tool.ComputerUse},
)

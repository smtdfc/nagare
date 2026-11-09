package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type GetAppInstanceInput struct {
	Name string `json:"name" jsonschema_description:"Exact name of the application to find and track (e.g. 'Calculator', 'Chrome', 'Firefox', 'Settings')."`
}

type GetAppInstanceOutput struct {
	InstanceID string `json:"instance_id" jsonschema_description:"Unique tracking identifier for this application instance to be used in other computer use tools."`
}

var GetAppInstanceTool = tool.DefineTool(
	"get_app_instance_tool",
	"Locate a running application by its name and obtain an instance ID for further UI automation and inspection.",
	func(ctx *context.ExecuteContext, args *GetAppInstanceInput, bindings tool.Bindings) (*GetAppInstanceOutput, error) {
		if args == nil || args.Name == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		appID, err := bindings.GetComputerUseOrchestrator().GetAppInstance(ctx, args.Name)
		if err != nil {
			return nil, err
		}

		return &GetAppInstanceOutput{InstanceID: appID}, nil
	},
	[]string{tool.ComputerUse},
)

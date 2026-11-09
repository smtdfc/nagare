package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type ClickElementInput struct {
	InstanceID string `json:"instance_id" jsonschema_description:"Unique tracking identifier of the target application instance obtained from get_app_instance_tool."`
	RefID      string `json:"ref_id" jsonschema_description:"Element reference ID obtained from get_ui_tree_tool or query_element_by_text_tool (e.g. '@e1', '@e2'). Do not invent arbitrary selectors."`
}

type ClickElementOutput struct {
	Success bool `json:"success" jsonschema_description:"Indicates whether the element was located and successfully clicked."`
}

var ClickElementTool = tool.DefineTool(
	"click_element_tool",
	"Click or invoke an interactive element inside an application using its reference ID (@e1, @e2...) obtained from get_ui_tree_tool or query_element_by_text_tool.",
	func(ctx *context.ExecuteContext, args *ClickElementInput, bindings tool.Bindings) (*ClickElementOutput, error) {
		if args == nil || args.InstanceID == "" || args.RefID == "" {
			return nil, fmt.Errorf("instance_id and ref_id cannot be empty")
		}

		err := bindings.GetComputerUseOrchestrator().ClickElement(ctx, args.InstanceID, args.RefID)
		if err != nil {
			return nil, err
		}

		return &ClickElementOutput{Success: true}, nil
	},
	[]string{tool.ComputerUse},
)

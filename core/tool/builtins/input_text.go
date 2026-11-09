package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type InputTextInput struct {
	InstanceID string `json:"instance_id,omitempty" jsonschema_description:"Optional application instance ID. If provided alongside ref_id, the text is targeted directly at that specific element."`
	RefID      string `json:"ref_id,omitempty" jsonschema_description:"Optional element reference ID obtained from get_ui_tree_tool or query_element_by_text_tool (e.g. '@e1', '@e2'). If omitted, text is typed into the currently active focus."`
	Text       string `json:"text" jsonschema_description:"Text content to enter or type into the target element or currently focused window."`
}

type InputTextOutput struct {
	Success bool `json:"success" jsonschema_description:"Indicates whether the text was successfully set or typed."`
}

var InputTextTool = tool.DefineTool(
	"input_text_tool",
	"Enter or type text into a specific application UI element (using its ref_id like '@e1') or directly into the currently focused control across the OS.",
	func(ctx *context.ExecuteContext, args *InputTextInput, bindings tool.Bindings) (*InputTextOutput, error) {
		if args == nil || args.Text == "" {
			return nil, fmt.Errorf("text cannot be empty")
		}

		err := bindings.GetComputerUseOrchestrator().InputText(ctx, args.InstanceID, args.RefID, args.Text)
		if err != nil {
			return nil, err
		}

		return &InputTextOutput{Success: true}, nil
	},
	[]string{tool.ComputerUse},
)

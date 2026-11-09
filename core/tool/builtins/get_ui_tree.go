package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/computer_use"
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type GetUITreeInput struct {
	InstanceID string `json:"instance_id" jsonschema_description:"Unique tracking identifier of the target application instance obtained from get_app_instance_tool."`
	MaxDepth   *int   `json:"max_depth,omitempty" jsonschema_description:"Optional maximum traversal depth for the accessibility tree (e.g. 3, 5, 10). Omit to dump the full tree."`
}

type GetUITreeOutput struct {
	Summary string                `json:"summary" jsonschema_description:"Indented plain-text summary of the simplified UI tree with assigned @e reference IDs (e.g. '@e1', '@e2') for easy element targeting."`
	Tree    *computer_use.UINode `json:"tree,omitempty" jsonschema_description:"Structured simplified hierarchy of UI elements with ref_id, role, name, value, bounds, and children."`
}

var GetUITreeTool = tool.DefineTool(
	"get_ui_tree_tool",
	"Inspect and dump a simplified accessibility UI tree of an application instance. Maps each element to an element reference ID (like '@e1', '@e2') and caches it for interaction tools (click_element_tool, input_text_tool).",
	func(ctx *context.ExecuteContext, args *GetUITreeInput, bindings tool.Bindings) (*GetUITreeOutput, error) {
		if args == nil || args.InstanceID == "" {
			return nil, fmt.Errorf("instance_id cannot be empty")
		}

		tree, summary, err := bindings.GetComputerUseOrchestrator().GetUITree(ctx, args.InstanceID, args.MaxDepth)
		if err != nil {
			return nil, err
		}

		return &GetUITreeOutput{
			Tree:    tree,
			Summary: summary,
		}, nil
	},
	[]string{tool.ComputerUse},
)

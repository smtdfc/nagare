package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/computer_use"
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type QueryElementByTextInput struct {
	InstanceID string `json:"instance_id" jsonschema_description:"Unique tracking identifier of the target application instance obtained from get_app_instance_tool."`
	Text       string `json:"text" jsonschema_description:"The text label, title, or value to search for within the application elements (e.g. 'File', 'Save', 'Submit', 'Search')."`
	Exact      *bool  `json:"exact,omitempty" jsonschema_description:"Whether to perform exact matching instead of substring matching (default is false, case-insensitive substring)."`
}

type QueryElementByTextOutput struct {
	Matches []computer_use.ElementMatch `json:"matches" jsonschema_description:"List of matched elements, each with a ref_id (e.g. '@e3') to use with click_element_tool or input_text_tool, plus role, name, value, and bounding box."`
	Total   int                         `json:"total" jsonschema_description:"Total number of matching elements found."`
}

var QueryElementByTextTool = tool.DefineTool(
	"query_element_by_text_tool",
	"Search for interactive UI elements (buttons, inputs, labels, menus) in an application by their visible text or name. Returns ready-to-use selectors and bounding boxes.",
	func(ctx *context.ExecuteContext, args *QueryElementByTextInput, bindings tool.Bindings) (*QueryElementByTextOutput, error) {
		if args == nil || args.InstanceID == "" || args.Text == "" {
			return nil, fmt.Errorf("instance_id and text cannot be empty")
		}

		matches, err := bindings.GetComputerUseOrchestrator().QueryElementByText(ctx, args.InstanceID, args.Text, args.Exact)
		if err != nil {
			return nil, err
		}

		return &QueryElementByTextOutput{
			Matches: matches,
			Total:   len(matches),
		}, nil
	},
	[]string{tool.ComputerUse},
)

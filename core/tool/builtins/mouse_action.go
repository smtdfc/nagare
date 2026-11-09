package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/computer_use"
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type MouseActionInput struct {
	Action       string   `json:"action" jsonschema_description:"Mouse operation to simulate. Allowed actions: 'move' (move pointer), 'click' (left-click), 'right_click' (context click), 'double_click' (fast double click), 'scroll' (scroll vertical/horizontal)."`
	X            *float64 `json:"x,omitempty" jsonschema_description:"Target X screen coordinate in pixels for move or click actions."`
	Y            *float64 `json:"y,omitempty" jsonschema_description:"Target Y screen coordinate in pixels for move or click actions."`
	ScrollDeltaX *float64 `json:"scroll_delta_x,omitempty" jsonschema_description:"Horizontal scroll delta in lines (positive scrolls right, negative scrolls left). Used when action is 'scroll'."`
	ScrollDeltaY *float64 `json:"scroll_delta_y,omitempty" jsonschema_description:"Vertical scroll delta in lines (positive scrolls down, negative scrolls up). Used when action is 'scroll'."`
}

type MouseActionOutput struct {
	Success bool `json:"success" jsonschema_description:"Indicates whether the mouse action was simulated successfully."`
}

var MouseActionTool = tool.DefineTool(
	"mouse_action_tool",
	"Simulate native OS mouse hardware events: cursor movement, left/right/double clicks at specific (x, y) coordinates, and scrolling.",
	func(ctx *context.ExecuteContext, args *MouseActionInput, bindings tool.Bindings) (*MouseActionOutput, error) {
		if args == nil || args.Action == "" {
			return nil, fmt.Errorf("action cannot be empty")
		}

		err := bindings.GetComputerUseOrchestrator().MouseAction(ctx, computer_use.MouseActionRequest{
			Action:       args.Action,
			X:            args.X,
			Y:            args.Y,
			ScrollDeltaX: args.ScrollDeltaX,
			ScrollDeltaY: args.ScrollDeltaY,
		})
		if err != nil {
			return nil, err
		}

		return &MouseActionOutput{Success: true}, nil
	},
	[]string{tool.ComputerUse},
)

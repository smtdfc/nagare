package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type GetActiveDisplayInput struct {
}

type GetActiveDisplayOutput struct {
	Num int `json:"num"`
}

var GetActiveDisplayTool = tool.DefineTool(
	"get_active_display_tool",
	"Get active display",
	func(ctx *context.ExecuteContext, args *GetActiveDisplayInput, bindings tool.Bindings) (*GetActiveDisplayOutput, error) {
		num, err := system.GetActiveDisplay()
		if err != nil {
			return nil, err
		}

		return &GetActiveDisplayOutput{Num: num}, nil
	},
	[]string{tool.SystemCategory},
)

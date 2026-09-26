package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type OpenBrowserInput struct {
	Url string `json:"url"`
}

type OpenBrowserOutput struct {
}

var OpenBrowserTool = tool.DefineTool(
	"open_browser_tool",
	"Open browser",
	func(ctx *context.ExecuteContext, args *OpenBrowserInput, _ tool.Bindings) (*OpenBrowserOutput, error) {
		err := system.OpenBrowser(args.Url)
		if err != nil {
			return nil, err
		}

		return &OpenBrowserOutput{}, nil
	},
)

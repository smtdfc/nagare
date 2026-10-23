package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type OpenBrowserInput struct {
	Urls []string `json:"urls"`
}

type OpenBrowserOutput struct {
}

var OpenBrowserTool = tool.DefineTool(
	"open_browser_tool",
	"Open browser",
	func(ctx *context.ExecuteContext, args *OpenBrowserInput, _ tool.Bindings) (*OpenBrowserOutput, error) {
		for _, url := range args.Urls {
			err := system.OpenBrowser(url)
			if err != nil {
				return nil, err
			}
		}

		return &OpenBrowserOutput{}, nil
	},
	[]string{tool.BrowserCategory},
)

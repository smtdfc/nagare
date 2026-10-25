package declarations

import (
	"fmt"
	"strings"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type OpenBrowserInput struct {
	Urls []string `json:"urls"`
}

type OpenBrowserOutput struct {
}

var OpenBrowserTool = tool.DefineTool(
	"open_browser_tool",
	"Open browser",
	func(ctx *context.ExecuteContext, args *OpenBrowserInput, bindings tool.Bindings) (*OpenBrowserOutput, error) {
		if args == nil || len(args.Urls) == 0 {
			return nil, fmt.Errorf("at least one URL is required")
		}

		for _, url := range args.Urls {
			if strings.TrimSpace(url) == "" {
				return nil, fmt.Errorf("URL cannot be empty")
			}
		}

		if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to open these URLs in the browser: %s", strings.Join(args.Urls, ", "))); err != nil {
			return nil, err
		}

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

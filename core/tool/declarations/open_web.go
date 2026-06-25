package declarations

import (
	"strings"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/cross-platform/system"
)

type OpenWebArgs struct {
	URL string `json:"url" jsonschema:"description=The target URL to be opened in the default browser"`
}

var open_web = tool.DeclareTool(
	"open_web",
	"Opens the specified URL in the default system web browser.",
	func(ctx domains.AgentContext, args OpenWebArgs) (any, error) {
		target := args.URL
		if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
			target = "https://" + target
		}

		err := system.OpenWeb(target)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"status": "success",
			"url":    target,
		}, nil
	},
	domains.STATIC_TOOL,
	domains.NO_GROUP,
)

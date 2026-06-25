package declarations

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/smtdfc/nagare/core/tool"
)

type OpenWebArgs struct {
	URL string `json:"url" jsonschema:"description=The target URL to be opened in the default browser"`
}

var open_web = tool.DeclareTool(
	"open_web",
	"Opens the specified URL in the default system web browser.",
	func(ctx context.Context, args OpenWebArgs) (any, error) {
		target := args.URL
		if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
			target = "https://" + target
		}

		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", target)
		case "darwin":
			cmd = exec.Command("open", target)
		case "linux":
			cmd = exec.Command("xdg-open", target)
		default:
			return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
		}

		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to open browser: %w", err)
		}

		return map[string]any{
			"status": "success",
			"url":    target,
		}, nil
	},
)

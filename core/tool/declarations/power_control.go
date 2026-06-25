package declarations

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/smtdfc/nagare/core/tool"
)

type PowerArgs struct {
	Action string `json:"action" jsonschema:"description=Action to perform: shutdown, reboot, sleep"`
}

func executePowerCommand(action string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		switch action {
		case "shutdown":
			cmd = exec.Command("shutdown", "/s", "/t", "0")
		case "reboot":
			cmd = exec.Command("shutdown", "/r", "/t", "0")
		case "sleep":
			cmd = exec.Command("powershell", "-Command", "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.Application]::SetSuspendState('Suspend', $false, $false)")
		}
	case "linux", "darwin":
		switch action {
		case "shutdown":
			cmd = exec.Command("shutdown", "now")
		case "reboot":
			cmd = exec.Command("reboot")
		case "sleep":
			cmd = exec.Command("systemctl", "suspend")
		}
	}

	if cmd == nil {
		return fmt.Errorf("unsupported action or OS")
	}

	return cmd.Run()
}

var power_control = tool.DeclareTool(
	"power_control",
	"Control system power states: shutdown, reboot, or sleep.",
	func(ctx context.Context, args PowerArgs) (any, error) {
		err := executePowerCommand(args.Action)
		if err != nil {
			return map[string]any{
				"action":  args.Action,
				"success": false,
				"error":   err.Error(),
			}, nil
		}

		return map[string]any{
			"action":  args.Action,
			"success": true,
		}, nil
	},
)

package declarations

import (
	"context"
	"fmt"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type KillProcessInput struct {
	PID int32 `json:"pid"` // PID of the process to be terminated
}

type KillProcessOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var KillProcessTool = tool.DefineTool(
	"kill_process_tool",
	"Terminate a running process by its PID",
	func(ctx context.Context, args *KillProcessInput) (*KillProcessOutput, error) {
		if args == nil || args.PID <= 0 {
			return nil, fmt.Errorf("invalid PID provided")
		}

		procMgr := system.NewProcessManager()
		name, err := procMgr.Kill(args.PID)
		if err != nil {
			return &KillProcessOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to terminate process PID %d: %v", args.PID, err),
			}, nil
		}

		return &KillProcessOutput{
			Success: true,
			Message: fmt.Sprintf("Successfully terminated process '%s' (PID: %d)", name, args.PID),
		}, nil
	},
)

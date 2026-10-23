package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type KillProcessInput struct {
	PID int32 `json:"pid" jsonschema_description:"The Process ID (PID) of the process to be terminated. This should be a positive integer representing the unique identifier of the running process."`
}

type KillProcessOutput struct {
	Success bool   `json:"success"`
	Message string `json:"messages"`
}

var KillProcessTool = tool.DefineTool(
	"kill_process_tool",
	"Terminate a running process by its PID",
	func(ctx *context.ExecuteContext, args *KillProcessInput, bindings tool.Bindings) (*KillProcessOutput, error) {
		if args == nil || args.PID <= 0 {
			return nil, fmt.Errorf("invalid PID provided")
		}

		procMgr := system.NewProcessManager()
		approved := bindings.RequestToolApproval(ctx, fmt.Sprintf("Allow the agent to terminate the process with PID: %d", args.PID))
		if !approved {
			return &KillProcessOutput{
				Success: false,
				Message: "User does not grant permission",
			}, nil
		}

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
	[]string{tool.ProcessManagementCategory},
)

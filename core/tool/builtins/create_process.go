package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type CreateProcessInput struct {
	Command    string   `json:"command"`
	Args       []string `json:"args,omitempty"`
	WorkingDir string   `json:"working_dir,omitempty"`
}

type CreateProcessOutput struct {
	PID     int32  `json:"pid"`
	Command string `json:"command"`
}

var CreateProcessTool = tool.DefineTool(
	"create_process_tool",
	"Start a detached process and return its PID. Pass the executable in command and arguments separately; do not use shell syntax.",
	func(ctx *context.ExecuteContext, args *CreateProcessInput, bindings tool.Bindings) (*CreateProcessOutput, error) {
		if args == nil || args.Command == "" {
			return nil, fmt.Errorf("command cannot be empty")
		}

		title := fmt.Sprintf("Allow the agent to start process %q", args.Command)
		if len(args.Args) > 0 {
			title = fmt.Sprintf("Allow the agent to start process %q with args %v", args.Command, args.Args)
		}
		if err := requireApproval(ctx, bindings, title); err != nil {
			return nil, err
		}

		pid, err := system.NewProcessManager().Start(args.Command, args.Args, args.WorkingDir)
		if err != nil {
			return nil, err
		}

		return &CreateProcessOutput{
			PID:     pid,
			Command: args.Command,
		}, nil
	},
	[]string{tool.ProcessManagementCategory},
)

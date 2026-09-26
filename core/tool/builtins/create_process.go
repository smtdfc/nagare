package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
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
	func(ctx *context.ExecuteContext, args *CreateProcessInput, _ tool.Bindings) (*CreateProcessOutput, error) {
		if args == nil || args.Command == "" {
			return nil, fmt.Errorf("command cannot be empty")
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
)

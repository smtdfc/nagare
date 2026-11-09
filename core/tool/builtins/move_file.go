package declarations

import (
	"fmt"
	"os"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type MoveFileInput struct {
	Source      string `json:"source" jsonschema_description:"Absolute or relative path of the file or directory to move or rename."`
	Destination string `json:"destination" jsonschema_description:"Absolute or relative destination path. If the destination is an existing directory the source is moved inside it; otherwise it is renamed to the given path."`
	Confirm     bool   `json:"confirm" jsonschema_description:"Must be true to allow the move or rename operation."`
}

type MoveFileOutput struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

var MoveFileTool = tool.DefineTool(
	"move_file_tool",
	"Move or rename a file or directory. Works for both files and directories. Confirm must be true.",
	func(ctx *context.ExecuteContext, args *MoveFileInput, bindings tool.Bindings) (*MoveFileOutput, error) {
		if args == nil || !args.Confirm {
			return nil, fmt.Errorf("move operation requires confirm=true")
		}
		if args.Source == "" || args.Destination == "" {
			return nil, fmt.Errorf("source and destination cannot be empty")
		}

		src, err := system.ResolveFilePath(args.Source)
		if err != nil {
			return nil, err
		}
		dst, err := system.ResolveFilePath(args.Destination)
		if err != nil {
			return nil, err
		}

		if _, err := os.Stat(src); err != nil {
			return nil, fmt.Errorf("source does not exist: %w", err)
		}

		if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to move %s → %s", src, dst)); err != nil {
			return nil, err
		}

		if err := os.Rename(src, dst); err != nil {
			return nil, fmt.Errorf("failed to move: %w", err)
		}

		return &MoveFileOutput{Source: src, Destination: dst}, nil
	},
	[]string{tool.FilesystemCategory},
)

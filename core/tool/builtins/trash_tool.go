package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type TrashToolInput struct {
	Path    string `json:"path" jsonschema_description:"Path of the file or directory to move into the system trash."`
	Confirm bool   `json:"confirm" jsonschema_description:"Must be true to allow the move to trash."`
}

type TrashToolOutput struct {
	Path   string `json:"path" jsonschema_description:"The original path that was trashed."`
	Target string `json:"target" jsonschema_description:"Where the path can be found after the move: a trash location path or 'Recycle Bin'."`
}

var TrashTool = tool.DefineTool(
	"trash_tool",
	"Move a file or directory into the operating system trash (Upload Trash spec on Linux, .Trash on macOS, Recycle Bin on Windows) instead of deleting it permanently. Confirm must be true.",
	func(ctx *context.ExecuteContext, args *TrashToolInput, bindings tool.Bindings) (*TrashToolOutput, error) {
		if args == nil || !args.Confirm {
			return nil, fmt.Errorf("moving to trash requires confirm=true")
		}
		if args.Path == "" {
			return nil, fmt.Errorf("path is required")
		}

		if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to move %s to the trash", args.Path)); err != nil {
			return nil, err
		}

		target, err := system.TrashMove(args.Path)
		if err != nil {
			return nil, err
		}

		return &TrashToolOutput{Path: args.Path, Target: target}, nil
	},
	[]string{tool.FilesystemCategory},
)

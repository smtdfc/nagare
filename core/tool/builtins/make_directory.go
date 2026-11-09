package declarations

import (
	"fmt"
	"os"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type MakeDirectoryInput struct {
	Path    string `json:"path" jsonschema_description:"Absolute or relative path of the directory to create. Intermediate parent directories are created automatically."`
	Confirm bool   `json:"confirm" jsonschema_description:"Must be true to allow the directory creation."`
}

type MakeDirectoryOutput struct {
	Path string `json:"path"`
}

var MakeDirectoryTool = tool.DefineTool(
	"make_directory_tool",
	"Create a directory (and any missing parent directories) at the given path. Confirm must be true.",
	func(ctx *context.ExecuteContext, args *MakeDirectoryInput, bindings tool.Bindings) (*MakeDirectoryOutput, error) {
		if args == nil || !args.Confirm {
			return nil, fmt.Errorf("directory creation requires confirm=true")
		}
		if args.Path == "" {
			return nil, fmt.Errorf("path cannot be empty")
		}

		path, err := system.ResolveFilePath(args.Path)
		if err != nil {
			return nil, err
		}

		if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to create directory %s", path)); err != nil {
			return nil, err
		}

		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}

		return &MakeDirectoryOutput{Path: path}, nil
	},
	[]string{tool.FilesystemCategory},
)

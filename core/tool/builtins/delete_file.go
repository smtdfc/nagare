package declarations

import (
	"fmt"
	"os"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type DeleteFileInput struct {
	Path    string `json:"path"`
	Confirm bool   `json:"confirm"`
}

type DeleteFileOutput struct {
	Path string `json:"path"`
}

var DeleteFileTool = tool.DefineTool(
	"delete_file_tool",
	"Delete a file at the provided path. Directories are not allowed and confirm must be true.",
	func(ctx *context.ExecuteContext, args *DeleteFileInput, _ tool.Bindings) (*DeleteFileOutput, error) {
		if args == nil || !args.Confirm {
			return nil, fmt.Errorf("file deletion requires confirm=true")
		}

		path, err := system.ResolveFilePath(args.Path)
		if err != nil {
			return nil, err
		}
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file: %w", err)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("directories cannot be deleted by this tool")
		}
		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("failed to delete file: %w", err)
		}

		return &DeleteFileOutput{Path: args.Path}, nil
	},
)

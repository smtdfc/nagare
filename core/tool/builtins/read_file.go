package declarations

import (
	"fmt"
	"os"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type ReadFileInput struct {
	Path string `json:"path"`
}

type ReadFileOutput struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Bytes   int    `json:"bytes"`
}

var ReadFileTool = tool.DefineTool(
	"read_file_tool",
	"Read a text file from the provided path.",
	func(ctx *context.ExecuteContext, args *ReadFileInput, _ tool.Bindings) (*ReadFileOutput, error) {
		if args == nil {
			return nil, fmt.Errorf("file path cannot be empty")
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
			return nil, fmt.Errorf("path is a directory")
		}
		if info.Size() > system.MaxFileToolBytes {
			return nil, fmt.Errorf("file exceeds %d byte limit", system.MaxFileToolBytes)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		return &ReadFileOutput{Path: args.Path, Content: string(content), Bytes: len(content)}, nil
	},
)

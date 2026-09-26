package declarations

import (
	"fmt"
	"os"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type WriteFileInput struct {
	Path      string `json:"path"`
	Content   string `json:"content"`
	Overwrite bool   `json:"overwrite"`
	Confirm   bool   `json:"confirm"`
}

type WriteFileOutput struct {
	Path  string `json:"path"`
	Bytes int    `json:"bytes"`
}

var WriteFileTool = tool.DefineTool(
	"write_file_tool",
	"Create or update a text file at the provided path. Confirm must be true.",
	func(ctx *context.ExecuteContext, args *WriteFileInput, _ tool.Bindings) (*WriteFileOutput, error) {
		if args == nil || !args.Confirm {
			return nil, fmt.Errorf("file write requires confirm=true")
		}
		if len(args.Content) > system.MaxFileToolBytes {
			return nil, fmt.Errorf("content exceeds %d byte limit", system.MaxFileToolBytes)
		}

		path, err := system.ResolveFilePath(args.Path)
		if err != nil {
			return nil, err
		}
		if !args.Overwrite {
			if _, err := os.Stat(path); err == nil {
				return nil, fmt.Errorf("file already exists; set overwrite=true")
			} else if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to check file: %w", err)
			}
		}

		if err := os.WriteFile(path, []byte(args.Content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write file: %w", err)
		}

		return &WriteFileOutput{Path: args.Path, Bytes: len(args.Content)}, nil
	},
)

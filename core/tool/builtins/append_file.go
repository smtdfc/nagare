package declarations

import (
	"fmt"
	"os"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type AppendFileInput struct {
	Path    string `json:"path" jsonschema_description:"Absolute or relative path of the file to append to. The file is created if it does not exist."`
	Content string `json:"content" jsonschema_description:"Text content to append at the end of the file."`
	Confirm bool   `json:"confirm" jsonschema_description:"Must be true to allow the append operation."`
}

type AppendFileOutput struct {
	Path  string `json:"path"`
	Bytes int    `json:"bytes"`
}

var AppendFileTool = tool.DefineTool(
	"append_file_tool",
	"Append text to the end of a file. Creates the file if it does not exist. Confirm must be true.",
	func(ctx *context.ExecuteContext, args *AppendFileInput, bindings tool.Bindings) (*AppendFileOutput, error) {
		if args == nil || !args.Confirm {
			return nil, fmt.Errorf("append operation requires confirm=true")
		}
		if args.Path == "" {
			return nil, fmt.Errorf("path cannot be empty")
		}
		if len(args.Content) > system.MaxFileToolBytes {
			return nil, fmt.Errorf("content exceeds %d byte limit", system.MaxFileToolBytes)
		}

		path, err := system.ResolveFilePath(args.Path)
		if err != nil {
			return nil, err
		}

		if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to append to %s", path)); err != nil {
			return nil, err
		}

		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open file for append: %w", err)
		}
		defer f.Close()

		n, err := f.WriteString(args.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to append to file: %w", err)
		}

		return &AppendFileOutput{Path: path, Bytes: n}, nil
	},
	[]string{tool.FilesystemCategory},
)

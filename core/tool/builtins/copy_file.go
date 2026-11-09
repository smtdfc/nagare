package declarations

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type CopyFileInput struct {
	Source      string `json:"source" jsonschema_description:"Absolute or relative path of the source file to copy."`
	Destination string `json:"destination" jsonschema_description:"Absolute or relative destination path for the new copy."`
	Overwrite   bool   `json:"overwrite" jsonschema_description:"Whether to overwrite the destination if it already exists."`
	Confirm     bool   `json:"confirm" jsonschema_description:"Must be true to allow the copy operation."`
}

type CopyFileOutput struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Bytes       int64  `json:"bytes"`
}

var CopyFileTool = tool.DefineTool(
	"copy_file_tool",
	"Copy a file to a new location. Only files are supported (not directories). Confirm must be true.",
	func(ctx *context.ExecuteContext, args *CopyFileInput, bindings tool.Bindings) (*CopyFileOutput, error) {
		if args == nil || !args.Confirm {
			return nil, fmt.Errorf("copy operation requires confirm=true")
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

		srcInfo, err := os.Stat(src)
		if err != nil {
			return nil, fmt.Errorf("source does not exist: %w", err)
		}
		if srcInfo.IsDir() {
			return nil, fmt.Errorf("source is a directory; only files can be copied")
		}
		if srcInfo.Size() > system.MaxFileToolBytes {
			return nil, fmt.Errorf("source file exceeds %d byte limit", system.MaxFileToolBytes)
		}

		if !args.Overwrite {
			if _, err := os.Stat(dst); err == nil {
				return nil, fmt.Errorf("destination already exists; set overwrite=true to replace it")
			}
		}

		if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to copy %s → %s", src, dst)); err != nil {
			return nil, err
		}

		// Ensure destination parent directory exists
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return nil, fmt.Errorf("failed to create destination directory: %w", err)
		}

		in, err := os.Open(src)
		if err != nil {
			return nil, fmt.Errorf("failed to open source: %w", err)
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			return nil, fmt.Errorf("failed to create destination: %w", err)
		}
		defer out.Close()

		n, err := io.Copy(out, in)
		if err != nil {
			return nil, fmt.Errorf("failed to copy file: %w", err)
		}

		return &CopyFileOutput{Source: src, Destination: dst, Bytes: n}, nil
	},
	[]string{tool.FilesystemCategory},
)

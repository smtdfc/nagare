package declarations

import (
	"fmt"
	"os"
	"time"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type StatFileInput struct {
	Path string `json:"path" jsonschema_description:"Absolute or relative path of the file or directory to inspect."`
}

type StatFileOutput struct {
	Path        string    `json:"path"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	IsDir       bool      `json:"is_dir"`
	Mode        string    `json:"mode"`
	ModifiedAt  time.Time `json:"modified_at"`
	Exists      bool      `json:"exists"`
}

var StatFileTool = tool.DefineTool(
	"stat_file_tool",
	"Get metadata of a file or directory: size, type, permissions, and modification time. Returns exists=false when the path does not exist.",
	func(ctx *context.ExecuteContext, args *StatFileInput, _ tool.Bindings) (*StatFileOutput, error) {
		if args == nil || args.Path == "" {
			return nil, fmt.Errorf("path cannot be empty")
		}

		path, err := system.ResolveFilePath(args.Path)
		if err != nil {
			return nil, err
		}

		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			return &StatFileOutput{Path: path, Exists: false}, nil
		}
		if err != nil {
			return nil, fmt.Errorf("failed to stat path: %w", err)
		}

		return &StatFileOutput{
			Path:       path,
			Name:       info.Name(),
			Size:       info.Size(),
			IsDir:      info.IsDir(),
			Mode:       info.Mode().String(),
			ModifiedAt: info.ModTime(),
			Exists:     true,
		}, nil
	},
	[]string{tool.FilesystemCategory},
)

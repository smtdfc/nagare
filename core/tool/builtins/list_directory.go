package declarations

import (
	"fmt"
	"os"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type ListDirectoryInput struct {
	Path string `json:"path,omitempty"`
}

type FileEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
}

type ListDirectoryOutput struct {
	Path    string      `json:"path"`
	Entries []FileEntry `json:"entries"`
}

var ListDirectoryTool = tool.DefineTool(
	"list_directory_tool",
	"List files and directories at the provided path.",
	func(ctx *context.ExecuteContext, args *ListDirectoryInput, _ tool.Bindings) (*ListDirectoryOutput, error) {
		path := ""
		if args != nil {
			path = args.Path
		}

		resolvedPath, err := os.Getwd()
		if path != "" {
			resolvedPath, err = system.ResolveFilePath(path)
		}
		if err != nil {
			return nil, err
		}

		entries, err := os.ReadDir(resolvedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to list directory: %w", err)
		}

		result := make([]FileEntry, 0, len(entries))
		for _, entry := range entries {
			result = append(result, FileEntry{Name: entry.Name(), IsDir: entry.IsDir()})
		}

		return &ListDirectoryOutput{Path: path, Entries: result}, nil
	},
)

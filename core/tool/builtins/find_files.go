package declarations

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type FindFilesInput struct {
	Directory string `json:"directory" jsonschema_description:"Root directory to search in. Defaults to the current working directory if omitted."`
	Pattern   string `json:"pattern" jsonschema_description:"Glob pattern to match file names (e.g. '*.go', '*.log', 'config.*'). Matched against the file name only, not the full path."`
	MaxDepth  int    `json:"max_depth,omitempty" jsonschema_description:"Maximum directory recursion depth. 0 or negative means unlimited. Defaults to unlimited."`
}

type FindFilesOutput struct {
	Directory string   `json:"directory"`
	Pattern   string   `json:"pattern"`
	Matches   []string `json:"matches"`
	Total     int      `json:"total"`
}

var FindFilesTool = tool.DefineTool(
	"find_files_tool",
	"Recursively search for files matching a glob pattern (e.g. '*.go', '*.log') inside a directory. Returns a list of matching absolute paths.",
	func(ctx *context.ExecuteContext, args *FindFilesInput, _ tool.Bindings) (*FindFilesOutput, error) {
		if args == nil || args.Pattern == "" {
			return nil, fmt.Errorf("pattern cannot be empty")
		}

		dir := args.Directory
		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("failed to get working directory: %w", err)
			}
		} else {
			resolved, err := system.ResolveFilePath(dir)
			if err != nil {
				return nil, err
			}
			dir = resolved
		}

		info, err := os.Stat(dir)
		if err != nil {
			return nil, fmt.Errorf("directory does not exist: %w", err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("path is not a directory")
		}

		var matches []string
		maxDepth := args.MaxDepth
		rootDepth := len(filepath.SplitList(filepath.ToSlash(dir)))

		err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // skip unreadable entries
			}
			if d.IsDir() {
				if maxDepth > 0 {
					currentDepth := len(filepath.SplitList(filepath.ToSlash(path))) - rootDepth
					if currentDepth >= maxDepth {
						return filepath.SkipDir
					}
				}
				return nil
			}
			matched, err := filepath.Match(args.Pattern, d.Name())
			if err != nil {
				return fmt.Errorf("invalid pattern: %w", err)
			}
			if matched {
				matches = append(matches, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory: %w", err)
		}

		return &FindFilesOutput{
			Directory: dir,
			Pattern:   args.Pattern,
			Matches:   matches,
			Total:     len(matches),
		}, nil
	},
	[]string{tool.FilesystemCategory},
)

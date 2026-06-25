package declarations

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool"
)

type OpenFileManagerArgs struct {
	Path string `json:"path" jsonschema:"description=Directory path to open, empty means home directory"`
}

func openFileManager(path string) error {
	if path == "" {
		var err error
		path, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	}

	switch runtime.GOOS {

	case "windows":
		// explorer <path>
		return exec.Command("explorer", path).Run()

	case "darwin":
		// open <path>
		return exec.Command("open", path).Run()

	case "linux":
		// xdg-open <path>
		return exec.Command("xdg-open", path).Run()

	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

var open_file_manager = tool.DeclareTool(
	"open_file_manager",
	"Open system file manager at specified path",
	func(ctx domains.AgentContext, args OpenFileManagerArgs) (any, error) {
		openFileManager(args.Path)
		return map[string]any{
			"status": "success",
			"path":   args.Path,
		}, nil
	},
	domains.DYNAMIC_TOOL,
	domains.ListCategory{domains.FILE_TOOL, domains.PC_TOOL},
)

package registry

import (
	"github.com/smtdfc/nagare/core/tool"
	declarations "github.com/smtdfc/nagare/core/tool/builtins"
)

type ToolItem struct {
	RequiresRouter bool
	Tool           tool.Tool
}

var Registry = map[string]ToolItem{}

func RegisterTool(tool tool.Tool, requiresRouter bool) {
	Registry[tool.GetName()] = ToolItem{
		RequiresRouter: requiresRouter,
		Tool:           tool,
	}
}

func init() {
	RegisterTool(declarations.FindTools, false)
	RegisterTool(declarations.ExecuteTool, false)

	// requires router
	RegisterTool(declarations.WeatherTool, true)
	RegisterTool(declarations.TimeTool, true)
	RegisterTool(declarations.ListProcessTool, true)
	RegisterTool(declarations.CreateProcessTool, true)
	RegisterTool(declarations.KillProcessTool, true)
	RegisterTool(declarations.PowerControlTool, true)
	RegisterTool(declarations.VolumeControlTool, true)
	RegisterTool(declarations.CreateTaskTool, true)
	RegisterTool(declarations.OpenBrowserTool, true)
	RegisterTool(declarations.ReadFileTool, true)
	RegisterTool(declarations.WriteFileTool, true)
	RegisterTool(declarations.ListDirectoryTool, true)
	RegisterTool(declarations.DeleteFileTool, true)
	RegisterTool(declarations.GetUserDirectoriesTool, true)
}

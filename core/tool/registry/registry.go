package registry

import (
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/tool"
	declarations "github.com/smtdfc/nagare/core/tool/builtins"
	"github.com/smtdfc/nagare/core/tool/plugin"
)

type ToolItem struct {
	RequiresRouter bool
	Tool           tool.Tool
	IsPluginTool   bool
	PluginID       string
}

var Registry = map[string]ToolItem{}

func RegisterTool(tool tool.Tool, requiresRouter bool) {
	Registry[tool.GetName()] = ToolItem{
		RequiresRouter: requiresRouter,
		Tool:           tool,
		IsPluginTool:   false,
		PluginID:       "",
	}
}

func RegisterPluginTool(name string, args string, description string, categories []string, pluginID string, requiresRouter bool, eventBus *event_bus.CoreEventBus) {
	Registry[name] = ToolItem{
		RequiresRouter: requiresRouter,
		Tool:           plugin.NewPluginTool(name, args, description, categories, pluginID, eventBus),
		IsPluginTool:   true,
		PluginID:       pluginID,
	}
}

func init() {
	RegisterTool(declarations.FindToolByCategories, false)
	RegisterTool(declarations.ExecuteTool, false)
	RegisterTool(declarations.CreateTaskTool, false)
	RegisterTool(declarations.TimeTool, false)

	// requires router
	RegisterTool(declarations.WeatherTool, true)
	RegisterTool(declarations.ListProcessTool, true)
	RegisterTool(declarations.CreateProcessTool, true)
	RegisterTool(declarations.KillProcessTool, true)
	RegisterTool(declarations.PowerControlTool, true)
	RegisterTool(declarations.VolumeControlTool, true)
	RegisterTool(declarations.OpenBrowserTool, true)
	RegisterTool(declarations.ReadFileTool, true)
	RegisterTool(declarations.WriteFileTool, true)
	RegisterTool(declarations.ListDirectoryTool, true)
	RegisterTool(declarations.DeleteFileTool, true)
	RegisterTool(declarations.GetUserDirectoriesTool, true)
}

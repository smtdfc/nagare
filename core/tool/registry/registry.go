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
	RegisterTool(declarations.CreateCronJobTool, false)
	RegisterTool(declarations.TimeTool, false)
	RegisterTool(declarations.SearchTool, false)
	RegisterTool(declarations.FetchContentTool, false)

	// requires router
	RegisterTool(declarations.WeatherTool, true)
	RegisterTool(declarations.ListProcessTool, true)
	RegisterTool(declarations.CreateProcessTool, true)
	RegisterTool(declarations.ExecuteShellCommandTool, true)
	RegisterTool(declarations.KillProcessTool, true)
	RegisterTool(declarations.PowerControlTool, true)
	RegisterTool(declarations.VolumeControlTool, true)
	RegisterTool(declarations.OpenBrowserTool, true)
	RegisterTool(declarations.ReadFileTool, true)
	RegisterTool(declarations.WriteFileTool, true)
	RegisterTool(declarations.AppendFileTool, true)
	RegisterTool(declarations.ListDirectoryTool, true)
	RegisterTool(declarations.DeleteFileTool, true)
	RegisterTool(declarations.MoveFileTool, true)
	RegisterTool(declarations.CopyFileTool, true)
	RegisterTool(declarations.MakeDirectoryTool, true)
	RegisterTool(declarations.StatFileTool, true)
	RegisterTool(declarations.FindFilesTool, true)
	RegisterTool(declarations.GetUserDirectoriesTool, true)
	RegisterTool(declarations.GetAppInstanceTool, true)
	RegisterTool(declarations.ListRunningAppsTool, true)
	RegisterTool(declarations.OpenAppTool, true)
	RegisterTool(declarations.GetUITreeTool, true)
	RegisterTool(declarations.ClickElementTool, true)
	RegisterTool(declarations.InputTextTool, true)
	RegisterTool(declarations.MouseActionTool, true)
	RegisterTool(declarations.KeyboardShortcutTool, true)
	RegisterTool(declarations.WindowControlTool, true)
	RegisterTool(declarations.QueryElementByTextTool, true)
	RegisterTool(declarations.HashTool, true)
	RegisterTool(declarations.ArchiveTool, true)
	RegisterTool(declarations.TrashTool, true)
	RegisterTool(declarations.DownloadFileTool, true)
	RegisterTool(declarations.HTTPRequestTool, true)
	RegisterTool(declarations.BrightnessControlTool, true)
	RegisterTool(declarations.GetActiveDisplayTool, true)
	RegisterTool(declarations.ScreenshotTool, true)
}

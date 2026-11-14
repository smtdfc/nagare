package registry

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

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

// registry is guarded by mu because plugin tools are registered and released at
// runtime while agents read the registry concurrently.
var (
	mu       sync.RWMutex
	registry = map[string]ToolItem{}
)

func RegisterTool(tool tool.Tool, requiresRouter bool) {
	mu.Lock()
	defer mu.Unlock()

	registry[tool.GetName()] = ToolItem{
		RequiresRouter: requiresRouter,
		Tool:           tool,
		IsPluginTool:   false,
		PluginID:       "",
	}
}

// RegisterPluginTool adds or replaces a tool owned by a plugin. A plugin may
// replace its own tool, but built-in tool names are reserved and another
// plugin's tool cannot be hijacked.
func RegisterPluginTool(name string, args string, description string, categories []string, pluginID string, requiresRouter bool, eventBus *event_bus.CoreEventBus) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("plugin tool name cannot be empty")
	}
	if strings.TrimSpace(pluginID) == "" {
		return errors.New("plugin id cannot be empty")
	}

	mu.Lock()
	defer mu.Unlock()

	if existing, exists := registry[name]; exists {
		if !existing.IsPluginTool {
			return fmt.Errorf("tool %q is reserved by a built-in tool", name)
		}
		if existing.PluginID != pluginID {
			return fmt.Errorf("tool %q is already registered by another plugin", name)
		}
	}

	registry[name] = ToolItem{
		RequiresRouter: requiresRouter,
		Tool:           plugin.NewPluginTool(name, args, description, categories, pluginID, eventBus),
		IsPluginTool:   true,
		PluginID:       pluginID,
	}

	return nil
}

// UnregisterPluginTools removes every tool owned by pluginID and returns the
// removed tool names.
func UnregisterPluginTools(pluginID string) []string {
	if strings.TrimSpace(pluginID) == "" {
		return nil
	}

	mu.Lock()
	defer mu.Unlock()

	removed := make([]string, 0)
	for name, item := range registry {
		if !item.IsPluginTool || item.PluginID != pluginID {
			continue
		}

		delete(registry, name)
		removed = append(removed, name)
	}

	sort.Strings(removed)
	return removed
}

func Get(name string) (ToolItem, bool) {
	mu.RLock()
	defer mu.RUnlock()

	item, exists := registry[name]
	return item, exists
}

// List returns a stable, name-sorted snapshot of every registered tool.
func List() []ToolItem {
	mu.RLock()
	defer mu.RUnlock()

	items := make([]ToolItem, 0, len(registry))
	for _, item := range registry {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Tool.GetName() < items[j].Tool.GetName()
	})

	return items
}

func Len() int {
	mu.RLock()
	defer mu.RUnlock()

	return len(registry)
}

func init() {
	RegisterTool(declarations.FindToolByCategories, false)
	RegisterTool(declarations.ExecuteTool, false)
	RegisterTool(declarations.CreateCronJobTool, false)
	RegisterTool(declarations.TimeTool, false)
	RegisterTool(declarations.SearchTool, false)
	RegisterTool(declarations.FetchContentTool, false)
	RegisterTool(declarations.GetOSInfoTool, false)

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
	RegisterTool(declarations.GetActiveDisplayTool, true)
	RegisterTool(declarations.ScreenshotTool, true)
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
}

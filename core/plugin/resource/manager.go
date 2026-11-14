package resource

import (
	"errors"
	"strings"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/prompt"
	tool_manager "github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/core/tool/registry"
	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

// ErrSystemPromptTooLarge is returned when a plugin injects a system prompt
// fragment that exceeds [contracts.MaxSystemPromptBytes].
var ErrSystemPromptTooLarge = errors.New("system prompt too large")

// Manager owns every runtime resource a plugin injects into the core: tools,
// tool categories and system prompt fragments. All plugin mutations go through
// this manager so a single Release call can revoke everything when a plugin
// stops, crashes or is uninstalled.
type Manager struct {
	toolMgr  *tool_manager.ToolManager
	eventBus *event_bus.CoreEventBus
	logger   *logger.BaseLogger
}

// @Injectable
func NewPluginResourceManager(
	toolMgr *tool_manager.ToolManager,
	eventBus *event_bus.CoreEventBus,
	logger *logger.BaseLogger,
) *Manager {
	return &Manager{
		toolMgr:  toolMgr,
		eventBus: eventBus,
		logger:   logger.With("module", "plugin-resource-manager"),
	}
}

// RegisterTool adds a plugin-provided tool to the dynamic tool registry.
func (m *Manager) RegisterTool(pluginID string, definition contracts.PluginTool) error {
	if err := registry.RegisterPluginTool(
		definition.Name,
		definition.Args,
		definition.Description,
		definition.Categories,
		pluginID,
		true,
		m.eventBus,
	); err != nil {
		return err
	}

	m.logger.Info("Registered plugin tool", "pluginID", pluginID, "name", definition.Name)
	return nil
}

// RegisterToolCategories adds the categories a plugin wants the model to route
// to. Built-in category descriptions always win.
func (m *Manager) RegisterToolCategories(pluginID string, categories map[string]string) {
	m.toolMgr.AddPluginToolCategories(pluginID, categories)
	m.logger.Info("Registered plugin tool categories", "pluginID", pluginID, "count", len(categories))
}

// SetSystemPrompt replaces the system prompt fragment injected by a plugin.
// An empty content clears the fragment.
func (m *Manager) SetSystemPrompt(pluginID, content string) error {
	if len(content) > contracts.MaxSystemPromptBytes {
		return ErrSystemPromptTooLarge
	}

	prompt.SetPluginSystemPrompt(pluginID, content)
	if strings.TrimSpace(content) == "" {
		m.logger.Info("Cleared plugin system prompt", "pluginID", pluginID)
		return nil
	}

	m.logger.Info("Injected plugin system prompt", "pluginID", pluginID, "bytes", len(content))
	return nil
}

// ClearSystemPrompt removes the system prompt fragment injected by a plugin.
func (m *Manager) ClearSystemPrompt(pluginID string) {
	prompt.ClearPluginSystemPrompt(pluginID)
	m.logger.Info("Cleared plugin system prompt", "pluginID", pluginID)
}

// Release revokes every resource owned by pluginID. It is safe to call for a
// plugin that never registered anything.
func (m *Manager) Release(pluginID string) {
	pluginID = strings.TrimSpace(pluginID)
	if pluginID == "" {
		return
	}

	removedTools := registry.UnregisterPluginTools(pluginID)
	m.toolMgr.RemovePluginToolCategories(pluginID)
	prompt.ClearPluginSystemPrompt(pluginID)

	if len(removedTools) > 0 {
		m.logger.Info(
			"Released plugin resources",
			"pluginID", pluginID,
			"tools", strings.Join(removedTools, ", "),
		)
		return
	}

	m.logger.Debug("Released plugin resources", "pluginID", pluginID)
}

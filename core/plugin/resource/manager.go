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

var ErrSystemPromptTooLarge = errors.New("system prompt too large")

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

func (m *Manager) RegisterToolCategories(pluginID string, categories map[string]string) {
	m.toolMgr.AddPluginToolCategories(pluginID, categories)
	m.logger.Info("Registered plugin tool categories", "pluginID", pluginID, "count", len(categories))
}

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

func (m *Manager) ClearSystemPrompt(pluginID string) {
	prompt.ClearPluginSystemPrompt(pluginID)
	m.logger.Info("Cleared plugin system prompt", "pluginID", pluginID)
}

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

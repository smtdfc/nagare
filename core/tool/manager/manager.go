package manager

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/smtdfc/nagare/core/computer_use"
	"github.com/smtdfc/nagare/core/context"
	cronjob_manager "github.com/smtdfc/nagare/core/cronjob/manager"
	discovery_manager "github.com/smtdfc/nagare/core/discovery/manager"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/hil"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/tool/registry"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

// experimentalNoteComputerUse flags every computer use tool as an
// experimental feature: GUI automation is slow, flaky, and costly for tasks
// that a native tool could perform directly. Keep the feature behind this
// warning until native "attached hands" tooling can serve those workflows.
const experimentalNoteComputerUse = "computer use tools are an experimental feature: they are slow, unreliable and token-heavy compared to native tools"

type ToolManager struct {
	mu                      sync.RWMutex
	cachedList              tool.ListTool
	logger                  *logger.BaseLogger
	cronJobMgr              *cronjob_manager.CronJobManager
	eventBus                *event_bus.CoreEventBus
	approvalMgr             *hil.Manager
	discoveryMgr            *discovery_manager.DiscoveryManager
	computerUseOrchestrator *computer_use.ComputerUseOrchestrator

	toolCategories       map[string]string
	pluginToolCategories map[string]map[string]string
}

func (t *ToolManager) createBindings() *ToolBindings {
	return &ToolBindings{
		cronJobMgr:              t.cronJobMgr,
		eventBus:                t.eventBus,
		toolMgr:                 t,
		approvalMgr:             t.approvalMgr,
		discoveryMgr:            t.discoveryMgr,
		computerUseOrchestrator: t.computerUseOrchestrator,
	}
}

func (t *ToolManager) GetListTool() tool.ListTool {
	t.mu.RLock()
	if t.cachedList != nil {
		t.mu.RUnlock()
		return t.cachedList
	}
	t.mu.RUnlock()

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.cachedList != nil {
		return t.cachedList
	}

	list := make(tool.ListTool, 0, len(registry.Registry))
	for _, item := range registry.Registry {
		if !item.RequiresRouter {
			list = append(list, item.Tool)
		}
	}

	t.cachedList = list
	return t.cachedList
}

func (t *ToolManager) Call(ctx *context.ExecuteContext, toolCall *tool.Call) *tool.Result {
	toolResultBuilder := tool.NewToolResultBuilder(toolCall.CallID, toolCall.Name)
	calledTool, isExist := registry.Registry[toolCall.Name]
	if !isExist {
		return toolResultBuilder.Failure(errors.New("tool doesn't exist")).Build()
	}

	if hasComputerUseCategory(calledTool.Tool.GetCategories()) {
		warnExperimentalComputerUse(t.logger, toolCall.Name)
	}

	bindings := t.createBindings()
	result, err := calledTool.Tool.WithBindings(bindings).Execute(ctx, toolCall.Args)
	if err != nil {
		return toolResultBuilder.Failure(err).Build()
	}

	return toolResultBuilder.Success(result).Build()
}

// hasComputerUseCategory reports whether any of the tool categories marks a
// computer use tool.
func hasComputerUseCategory(categories []string) bool {
	for _, category := range categories {
		if category == tool.ComputerUse {
			return true
		}
	}
	return false
}

func (t *ToolManager) FindToolsByKeywords(keywords []string) ([]tool.Metadata, error) {
	list := make([]tool.Metadata, 0)

	for _, item := range registry.Registry {
		if !item.RequiresRouter {
			continue
		}

		nameMatch := helpers.ContainsAnyWord(item.Tool.GetName(), keywords)
		descMatch := helpers.ContainsAnyWord(item.Tool.GetDescription(), keywords)
		categoriesMatch := false
		for _, category := range item.Tool.GetCategories() {
			if helpers.ContainsAnyWord(category, keywords) {
				categoriesMatch = true
				break
			}
		}

		if nameMatch || descMatch || categoriesMatch {
			list = append(list, tool.Metadata{
				Name:         item.Tool.GetName(),
				Description:  item.Tool.GetDescription(),
				Args:         item.Tool.GetArgsSchema(),
				PluginID:     item.PluginID,
				IsPluginTool: item.IsPluginTool,
			})
		}
	}

	return list, nil
}

func (t *ToolManager) GetCategoriesPrompt() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var prompt strings.Builder
	for name, desc := range t.toolCategories {
		prompt.WriteString(fmt.Sprintf("%s (%s), ", name, desc))
	}

	return prompt.String()
}

func (t *ToolManager) AddCategory(category string, description string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.toolCategories[category] = description
}

func (t *ToolManager) AddPluginToolCategories(pluginID string, categories map[string]string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.pluginToolCategories == nil {
		t.pluginToolCategories = make(map[string]map[string]string)
	}

	t.pluginToolCategories[pluginID] = categories
}

// @Injectable
func NewToolManager(logger *logger.BaseLogger, cronJobMgr *cronjob_manager.CronJobManager, eventBus *event_bus.CoreEventBus, approvalMgr *hil.Manager, discoveryMgr *discovery_manager.DiscoveryManager, computerUseOrchestrator *computer_use.ComputerUseOrchestrator) *ToolManager {
	return &ToolManager{
		logger:                  logger.With("module", "tool-manager"),
		cronJobMgr:              cronJobMgr,
		eventBus:                eventBus,
		approvalMgr:             approvalMgr,
		discoveryMgr:            discoveryMgr,
		computerUseOrchestrator: computerUseOrchestrator,
		toolCategories: map[string]string{
			tool.ProcessManagementCategory: "Process Management",
			tool.PowerManagementCategory:   "Power Management",
			tool.AudioManagementCategory:   "Audio Management",
			tool.FilesystemCategory:        "Filesystem",
			tool.BrowserCategory:           "Browser",
			tool.WeatherCategory:           "Weather",
			tool.NetworkingCategory:        "Networking",
			tool.CronJobManagementCategory: "CronJob Management",
			tool.RoutingCategory:           "Tool Routing",
			tool.TimingCategory:            "Timing",
			tool.SystemCategory:            "System",
			tool.ResearchCategory:          "Research",
			tool.ComputerUse:               "Computer use",
		},
		pluginToolCategories: make(map[string]map[string]string),
	}
}

var experimentalComputerUseWarnOnce sync.Once

func warnExperimentalComputerUse(callLogger *logger.BaseLogger, toolName string) {
	experimentalComputerUseWarnOnce.Do(func() {
		var names []string
		for name, item := range registry.Registry {
			if hasComputerUseCategory(item.Tool.GetCategories()) {
				names = append(names, name)
			}
		}

		callLogger.Warn("Experimental tools registered", "category", tool.ComputerUse, "tools", strings.Join(names, ", "), "note", experimentalNoteComputerUse)
	})

	callLogger.Warn("Experimental feature used", "tool", toolName, "category", tool.ComputerUse, "note", experimentalNoteComputerUse)
}

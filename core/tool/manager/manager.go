package manager

import (
	context2 "context"
	"errors"
	"strings"
	"sync"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	task_manager "github.com/smtdfc/nagare/core/task/manager"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/tool/registry"
	"github.com/smtdfc/nagare/shared/helpers"
)

type ToolManager struct {
	mu         sync.RWMutex
	cachedList tool.ListTool
	logger     *logger.BaseLogger
	taskMgr    *task_manager.TaskManager
	eventBus   *event_bus.CoreEventBus

	toolCategories       []string
	pluginToolCategories map[string][]string
}

func (t *ToolManager) createBindings() *ToolBindings {
	return &ToolBindings{
		taskMgr:  t.taskMgr,
		eventBus: t.eventBus,
		toolMgr:  t,
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

func (t *ToolManager) Call(ctx *context.ExecuteContext, toolCall *tool.ToolCall) *tool.Result {
	toolResultBuilder := tool.NewToolResultBuilder(toolCall.CallID, toolCall.Name)
	calledTool, isExist := registry.Registry[toolCall.Name]
	if !isExist {
		return toolResultBuilder.Failure(errors.New("tool doesn't exist")).Build()
	}

	bindings := t.createBindings()
	result, err := calledTool.Tool.WithBindings(bindings).Execute(ctx, toolCall.Args)
	if err != nil {
		return toolResultBuilder.Failure(err).Build()
	}

	return toolResultBuilder.Success(result).Build()
}

func (t *ToolManager) FindToolsByKeywords(ctx context2.Context, keywords []string) ([]tool.ToolMetadata, error) {
	list := make([]tool.ToolMetadata, 0)

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
			list = append(list, tool.ToolMetadata{
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

func (t *ToolManager) GetCategoriesString() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	allCategories := make([]string, len(t.toolCategories))
	copy(allCategories, t.toolCategories)

	for _, pluginCategories := range t.pluginToolCategories {
		allCategories = append(allCategories, pluginCategories...)
	}

	categorySet := make(map[string]struct{})
	for _, category := range allCategories {
		categorySet[category] = struct{}{}
	}

	uniqueCategories := make([]string, 0, len(categorySet))
	for category := range categorySet {
		uniqueCategories = append(uniqueCategories, category)
	}

	return strings.Join(uniqueCategories, ", ")
}

func (t *ToolManager) GetPluginToolCategoriesString(pluginID string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if categories, exists := t.pluginToolCategories[pluginID]; exists {
		return strings.Join(categories, ", ")
	}
	return strings.Join(t.toolCategories, ", ")
}

func (t *ToolManager) AddCategory(category string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, existingCategory := range t.toolCategories {
		if existingCategory == category {
			return
		}
	}

	t.toolCategories = append(t.toolCategories, category)
}

func (t *ToolManager) AddPluginToolCategories(pluginID string, categories []string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.pluginToolCategories[pluginID]; exists {
		t.pluginToolCategories[pluginID] = append(t.pluginToolCategories[pluginID], categories...)
	} else {
		t.pluginToolCategories[pluginID] = categories
	}
}

// @Injectable
func NewToolManager(logger *logger.BaseLogger, taskMgr *task_manager.TaskManager, eventBus *event_bus.CoreEventBus) *ToolManager {
	return &ToolManager{
		logger:   logger.With("module", "tool-manager"),
		taskMgr:  taskMgr,
		eventBus: eventBus,
		toolCategories: []string{
			tool.ProcessManagementCategory,
			tool.PowerManagementCategory,
			tool.AudioManagementCategory,
			tool.FilesystemCategory,
			tool.BrowserCategory,
			tool.WeatherCategory,
			tool.NetworkingCategory,
			tool.TaskManagementCategory,
			tool.ToolRoutingCategory,
			tool.TimingCategory,
		},
		pluginToolCategories: make(map[string][]string),
	}
}

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
	"github.com/smtdfc/nagare/pkgs/helpers"
)

type ToolManager struct {
	mu         sync.RWMutex
	cachedList tool.ListTool
	logger     *logger.BaseLogger
	taskMgr    *task_manager.TaskManager
	eventBus   *event_bus.CoreEventBus

	toolCategories       map[string]string
	pluginToolCategories map[string]map[string]string
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

func (t *ToolManager) GetCategoriesPrompt() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var prompt strings.Builder
	for category, description := range t.toolCategories {
		prompt.WriteString("<category name=\"")
		prompt.WriteString(category)
		prompt.WriteString("\">")
		prompt.WriteString(description)
		prompt.WriteString("</category>")
	}

	for pluginID, categories := range t.pluginToolCategories {
		for categoryName, categoryDescription := range categories {
			prompt.WriteString("<category name=\"")
			prompt.WriteString(categoryName)
			prompt.WriteString("\" plugin_id=\"")
			prompt.WriteString(pluginID)
			prompt.WriteString("\">")
			prompt.WriteString(categoryDescription)
			prompt.WriteString("</category>")
		}
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
func NewToolManager(logger *logger.BaseLogger, taskMgr *task_manager.TaskManager, eventBus *event_bus.CoreEventBus) *ToolManager {
	return &ToolManager{
		logger:   logger.With("module", "tool-manager"),
		taskMgr:  taskMgr,
		eventBus: eventBus,
		toolCategories: map[string]string{
			tool.ProcessManagementCategory: "Process Management",
			tool.PowerManagementCategory:   "Power Management",
			tool.AudioManagementCategory:   "Audio Management",
			tool.FilesystemCategory:        "Filesystem",
			tool.BrowserCategory:           "Browser",
			tool.WeatherCategory:           "Weather",
			tool.NetworkingCategory:        "Networking",
			tool.TaskManagementCategory:    "Task Management",
			tool.ToolRoutingCategory:       "Tool Routing",
			tool.TimingCategory:            "Timing",
		},
		pluginToolCategories: make(map[string]map[string]string),
	}
}

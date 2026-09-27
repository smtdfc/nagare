package manager

import (
	context2 "context"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	task2 "github.com/smtdfc/nagare/core/task"
	task_manager "github.com/smtdfc/nagare/core/task/manager"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/tool/registry"
	"github.com/smtdfc/nagare/shared/helpers"
)

type ToolBindings struct {
	taskMgr  *task_manager.TaskManager
	toolMgr  *ToolManager
	eventBus *event_bus.CoreEventBus
}

func (t *ToolBindings) FindToolsByKeywords(ctx *context.ExecuteContext, keywords []string) ([]tool.ToolMetadata, error) {
	return t.toolMgr.FindToolsByKeywords(ctx.Context, keywords)
}

func (t ToolBindings) RefreshTask(ctx *context.ExecuteContext) {
	t.eventBus.Publish(ctx, event_bus.RefreshTaskEvent, &event_bus.RefreshTaskEventPayload{})
}

func (t ToolBindings) CreateTask(ctx *context.ExecuteContext, sessionID, name, prompt string, triggerBy string, repeat bool, repeatRule string, startTime string, endTime string) (string, error) {
	var triggerSources = []string{"scheduled"}
	var repeatRules = []string{"no_repeat", "daily"}
	if !slices.Contains(triggerSources, triggerBy) {
		return "", errors.New("trigger source does not valid")
	}

	r := task2.NoRepeat
	if repeat {
		if !slices.Contains(repeatRules, repeatRule) {
			return "", errors.New("repeat rule does not valid")
		}
		r = task2.MapTaskRepeatRuleToTaskRepeatRule(repeatRule)
	}

	var parsedStartTime *time.Time
	if startTime != "" {
		tStart, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			tStart, err = time.Parse("2006-01-02 15:04:05", startTime)
			if err != nil {
				return "", errors.New("startTime format is invalid, expected RFC3339 or YYYY-MM-DD HH:MM:SS")
			}
		}
		parsedStartTime = &tStart
	}

	var parsedEndTime *time.Time
	if endTime != "" {
		tEnd, err := time.Parse(time.RFC3339, endTime)
		if err != nil {
			tEnd, err = time.Parse("2006-01-02 15:04:05", endTime)
			if err != nil {
				return "", errors.New("endTime format is invalid, expected RFC3339 or YYYY-MM-DD HH:MM:SS")
			}
		}
		parsedEndTime = &tEnd
	}

	task, err := t.taskMgr.Create(ctx, name, sessionID, prompt, &task2.TaskTriggerRule{
		By:        task2.MapTaskTriggerSourceToTaskSource(triggerBy),
		StartTime: parsedStartTime,
		EndTime:   parsedEndTime,
		Repeat:    r,
	})
	if err != nil {
		return "", err
	}

	return task.ID.String(), nil
}

func (t ToolBindings) GetTaskManager() *task_manager.TaskManager {
	return t.taskMgr
}

func (t ToolBindings) CallTool(ctx *context.ExecuteContext, toolName string, args string) *tool.Result {
	return t.toolMgr.Call(ctx, &tool.ToolCall{
		CallID: time.Now().Format("20060102150405"),
		Name:   toolName,
		Args:   args,
	})
}

type ToolManager struct {
	mu                sync.RWMutex
	cachedList        tool.ListTool
	logger            *logger.BaseLogger
	taskMgr           *task_manager.TaskManager
	eventBus          *event_bus.CoreEventBus
	dynamicToolRepo   *repositories.DynamicToolRepository
	dynamicToolMapper *mappers.DynamicToolMapper
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

func (t *ToolManager) UseDynamicTools(ctx context2.Context, dynamicTools []*tool.DynamicTool) error {
	err := t.dynamicToolRepo.UpsertBatch(t.dynamicToolMapper.ToEntities(dynamicTools))
	if err != nil {
		return custom_errors.ErrAddDynamicToolFailed
	}

	return nil
}

func (t *ToolManager) FindToolsByKeywords(ctx context2.Context, keywords []string) ([]tool.ToolMetadata, error) {
	list := make([]tool.ToolMetadata, 0)

	for _, item := range registry.Registry {
		if !item.RequiresRouter {
			continue
		}

		nameMatch := helpers.ContainsAnyKeyword(item.Tool.GetName(), keywords)
		descMatch := helpers.ContainsAnyKeyword(item.Tool.GetDescription(), keywords)

		if nameMatch || descMatch {
			list = append(list, tool.ToolMetadata{
				Name:        item.Tool.GetName(),
				Description: item.Tool.GetDescription(),
				Args:        item.Tool.GetArgsSchema(),
			})
		}
	}

	return list, nil
}

// @Injectable
func NewToolManager(logger *logger.BaseLogger, taskMgr *task_manager.TaskManager, eventBus *event_bus.CoreEventBus, dynamicToolRepo *repositories.DynamicToolRepository, dynamicToolMapper *mappers.DynamicToolMapper) *ToolManager {
	return &ToolManager{
		logger:            logger.With("module", "tool-manager"),
		taskMgr:           taskMgr,
		eventBus:          eventBus,
		dynamicToolRepo:   dynamicToolRepo,
		dynamicToolMapper: dynamicToolMapper,
	}
}

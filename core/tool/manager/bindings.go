package manager

import (
	"errors"
	"slices"
	"time"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/task"
	task_manager "github.com/smtdfc/nagare/core/task/manager"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/tool/registry"
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

	r := task.NoRepeat
	if repeat {
		if !slices.Contains(repeatRules, repeatRule) {
			return "", errors.New("repeat rule does not valid")
		}
		r = task.MapTaskRepeatRuleToTaskRepeatRule(repeatRule)
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

	task, err := t.taskMgr.Create(ctx, name, sessionID, prompt, &task.TaskTriggerRule{
		By:        task.MapTaskTriggerSourceToTaskSource(triggerBy),
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

func (t ToolBindings) FindToolsByCategories(ctx *context.ExecuteContext, categories []string) ([]tool.ToolMetadata, error) {
	list := make([]tool.ToolMetadata, 0)
	for _, item := range registry.Registry {
		if !item.RequiresRouter {
			continue
		}

		categoriesMatch := false
		for _, category := range item.Tool.GetCategories() {
			if slices.Contains(categories, category) {
				categoriesMatch = true
				break
			}
		}

		if categoriesMatch {
			list = append(list, tool.ToolMetadata{
				Name:        item.Tool.GetName(),
				Description: item.Tool.GetDescription(),
				Args:        item.Tool.GetArgsSchema(),
			})
		}
	}
	return list, nil
}

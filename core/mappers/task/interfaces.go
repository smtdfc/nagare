package task

import (
	"github.com/smtdfc/nagare/core/model_provider/provider"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/task"
)

func parseTaskStatus(s string) task.Status {
	return task.MapStringToTaskStatus(s)
}

func formatTaskStatus(s task.Status) string {
	return s.ToString()
}

func formatTaskSource(s *task.TriggerSource) string {
	if s == nil {
		return ""
	}
	return s.ToString()
}

func formatRepeatRule(r *task.RepeatRule) string {
	if r == nil {
		return ""
	}
	return r.ToString()
}

func mapEntityToTriggerRule(entity *entities.Task) *task.TriggerRule {
	if entity == nil {
		return nil
	}
	return &task.TriggerRule{
		By:        task.MapTaskTriggerSourceToTaskSource(entity.TriggerBy),
		StartTime: entity.StartTime,
		EndTime:   entity.EndTime,
		Repeat:    task.MapTaskRepeatRuleToTaskRepeatRule(entity.RepeatType),
		EventName: entity.TriggerByEventName,
	}
}

func parseCompatible(s string) provider.ModelProviderCompatible {
	return provider.GetCompatibleFromString(s)
}

// goverter:converter
// goverter:output:file ./mapper.generated.go
// goverter:ignoreMissing
// goverter:skipCopySameType
// goverter:useZeroValueOnPointerInconsistency
type ITaskMapper interface {
	// goverter:map Status | parseTaskStatus
	// goverter:map . TriggerRule | mapEntityToTriggerRule
	ToDomain(entity *entities.Task) *task.Task

	// goverter:map Status | formatTaskStatus
	// goverter:map TriggerRule.By TriggerBy | formatTaskSource
	// goverter:map TriggerRule.Repeat RepeatType | formatRepeatRule
	// goverter:map TriggerRule.EventName TriggerByEventName
	// goverter:map TriggerRule.StartTime StartTime
	// goverter:map TriggerRule.EndTime EndTime
	ToEntity(domain *task.Task) *entities.Task

	ToDomains(entities []*entities.Task) []*task.Task
	ToEntities(domains []*task.Task) []*entities.Task
}

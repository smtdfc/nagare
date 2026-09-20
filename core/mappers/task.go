package mappers

import (
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/task"
	"github.com/smtdfc/nagare/shared/helpers"
)

type TaskMapper struct {
}

func (t *TaskMapper) ToDomain(entity *entities.Task) *task.Task {
	return &task.Task{
		ID:        entity.ID,
		Name:      entity.Name,
		IsActive:  entity.IsActive,
		Prompt:    entity.Prompt,
		SessionID: entity.SessionID,
		Status:    task.MapStringToTaskStatus(entity.Status),
		TriggerRule: &task.TaskTriggerRule{
			By:        task.MapTaskTriggerSourceToTaskSource(entity.TriggerBy),
			StartTime: entity.StartTime,
			EndTime:   entity.EndTime,
			Repeat:    task.MapTaskRepeatRuleToTaskRepeatRule(entity.RepeatType),
			EventName: entity.TriggerByEventName,
		},
		NextRunTime: entity.NextRunTime,
	}
}

func (t *TaskMapper) ToEntity(domain *task.Task) *entities.Task {
	entity := &entities.Task{
		ID:        domain.ID,
		Name:      domain.Name,
		IsActive:  domain.IsActive,
		Prompt:    domain.Prompt,
		SessionID: domain.SessionID,
		Status:    domain.Status.ToString(),
	}

	if domain.TriggerRule != nil {
		entity.RepeatType = domain.TriggerRule.Repeat.ToString()
		entity.StartTime = domain.TriggerRule.StartTime
		entity.EndTime = domain.TriggerRule.EndTime
		entity.TriggerBy = domain.TriggerRule.By.ToString()
		entity.TriggerByEventName = domain.TriggerRule.EventName
		entity.NextRunTime = domain.NextRunTime
	}

	return entity
}

func (t *TaskMapper) ToDomains(entities []*entities.Task) []*task.Task {
	return helpers.SliceMap(entities, t.ToDomain)
}

func (t *TaskMapper) ToEntities(entities []*task.Task) []*entities.Task {
	return helpers.SliceMap(entities, t.ToEntity)
}

// @Injectable
func NewTaskMapper() *TaskMapper {
	return &TaskMapper{}
}

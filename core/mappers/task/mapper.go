//go:build !goverter

package task

import (
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/task"
)

type TaskMapper struct {
	converter TaskConverter
}

// @Injectable
func NewTaskMapper() *TaskMapper {
	return &TaskMapper{
		converter: &TaskConverterImpl{},
	}
}

func (t *TaskMapper) ToDomain(entity *entities.Task) *task.Task {
	return t.converter.ToDomain(entity)
}

func (t *TaskMapper) ToEntity(domain *task.Task) *entities.Task {
	return t.converter.ToEntity(domain)
}

func (t *TaskMapper) ToDomains(entities []*entities.Task) []*task.Task {
	return t.converter.ToDomains(entities)
}

func (t *TaskMapper) ToEntities(entities []*task.Task) []*entities.Task {
	return t.converter.ToEntities(entities)
}

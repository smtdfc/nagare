package task

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID
	Name      string
	SessionID uuid.UUID
	IsActive  bool

	Prompt      string
	Status      TaskStatus
	TriggerRule *TaskTriggerRule

	NextRunTime *time.Time
}

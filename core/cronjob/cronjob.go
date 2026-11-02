package cronjob

import (
	"time"

	"github.com/google/uuid"
)

type CronJob struct {
	ID          uuid.UUID
	Name        string
	SessionID   uuid.UUID
	IsActive    bool
	Prompt      string
	Status      Status
	CronExpr    string
	NextRunTime *time.Time
	LastRunTime *time.Time
}

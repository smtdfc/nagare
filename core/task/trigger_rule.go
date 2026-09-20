package task

import "time"

type TaskTriggerRule struct {
	By        TaskTriggerSource
	EventName string
	StartTime *time.Time
	EndTime   *time.Time
	Repeat    TaskRepeatRule
}

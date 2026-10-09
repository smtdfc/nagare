package task

import "time"

type TriggerRule struct {
	By        TriggerSource
	EventName string
	StartTime *time.Time
	EndTime   *time.Time
	Repeat    RepeatRule
}

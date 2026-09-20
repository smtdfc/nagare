package task

type TaskTriggerSource string

const (
	Scheduled TaskTriggerSource = "scheduled"
	Event     TaskTriggerSource = "event"
)

func (t TaskTriggerSource) ToString() string {
	return string(t)
}

func MapTaskTriggerSourceToTaskSource(s string) TaskTriggerSource {
	switch s {
	case "scheduled":
		return Scheduled
	case "event":
		return Event
	default:
		return Scheduled
	}
}

package task

type TriggerSource string

const (
	Scheduled TriggerSource = "scheduled"
	Event     TriggerSource = "event"
)

func (t TriggerSource) ToString() string {
	return string(t)
}

func MapTaskTriggerSourceToTaskSource(s string) TriggerSource {
	switch s {
	case "scheduled":
		return Scheduled
	case "event":
		return Event
	default:
		return Scheduled
	}
}

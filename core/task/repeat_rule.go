package task

type TaskRepeatRule string

const (
	NoRepeat TaskRepeatRule = "no_repeat"
	Daily    TaskRepeatRule = "daily"
	Hourly   TaskRepeatRule = "hourly"
)

func (t TaskRepeatRule) ToString() string {
	return string(t)
}

func MapTaskRepeatRuleToTaskRepeatRule(s string) TaskRepeatRule {
	switch s {
	case "no_repeat":
		return NoRepeat
	case "daily":
		return Daily
	case "hourly":
		return Hourly
	default:
		return NoRepeat
	}
}

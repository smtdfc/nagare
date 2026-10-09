package task

type RepeatRule string

const (
	NoRepeat RepeatRule = "no_repeat"
	Daily    RepeatRule = "daily"
	Hourly   RepeatRule = "hourly"
)

func (t RepeatRule) ToString() string {
	return string(t)
}

func MapTaskRepeatRuleToTaskRepeatRule(s string) RepeatRule {
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

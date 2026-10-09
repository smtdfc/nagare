package task

type Status string

const (
	Pending Status = "pending"
	Running Status = "running"
	Queued  Status = "queued"
)

func (t Status) ToString() string {
	return string(t)
}

func MapStringToTaskStatus(s string) Status {
	switch s {
	case "pending":
		return Pending
	case "running":
		return Running
	default:
		return Pending
	}
}

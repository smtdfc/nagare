package task

type TaskStatus string

const (
	Pending TaskStatus = "pending"
	Running TaskStatus = "running"
	Queued  TaskStatus = "queued"
)

func (t TaskStatus) ToString() string {
	return string(t)
}

func MapStringToTaskStatus(s string) TaskStatus {
	switch s {
	case "pending":
		return Pending
	case "running":
		return Running
	default:
		return Pending
	}
}

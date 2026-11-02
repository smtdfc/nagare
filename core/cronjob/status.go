package cronjob

type Status string

const (
	Pending Status = "pending"
	Running Status = "running"
	Queued  Status = "queued"
)

func (s Status) ToString() string {
	return string(s)
}

func MapStringToCronJobStatus(s string) Status {
	switch s {
	case "pending":
		return Pending
	case "running":
		return Running
	case "queued":
		return Queued
	default:
		return Pending
	}
}

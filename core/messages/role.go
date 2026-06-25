package messages

type Role string

const (
	USER      Role = "user"
	AGENT     Role = "agent"
	SYSTEM    Role = "system"
	DEVELOPER Role = "developer"
)

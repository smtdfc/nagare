package messages

type Role int

func (r Role) String() string {
	switch r {
	case AGENT:
		return "Agent"
	case USER:
		return "User"
	case SYSTEM:
		return "System"

	default:
		return "Unknown"
	}
}

const (
	AGENT Role = iota
	USER
	SYSTEM
	DEVELOPER
)

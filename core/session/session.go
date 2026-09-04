package session

type SessionOwnerType string

const (
	USER    SessionOwnerType = "user"
	PLUGIN  SessionOwnerType = "plugin"
	SYSTEM  SessionOwnerType = "system"
	UNKNOWN SessionOwnerType = "unknown"
)

func (t SessionOwnerType) ToString() string {
	return string(t)
}

func GetOwnerType(raw string) SessionOwnerType {
	switch raw {
	case "user":
		return USER
	case "plugin":
		return PLUGIN
	case "system":
		return SYSTEM
	default:
		return UNKNOWN
	}
}

type SessionInfo struct {
	ID        string
	Title     string
	OwnerID   string
	OwnerType SessionOwnerType
	ChannelID string
	IsArchive bool
}

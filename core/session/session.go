package session

import "github.com/google/uuid"

type OwnerType string

const (
	USER    OwnerType = "user"
	PLUGIN  OwnerType = "plugin"
	SYSTEM  OwnerType = "system"
	UNKNOWN OwnerType = "unknown"
)

func (t OwnerType) ToString() string {
	return string(t)
}

func GetOwnerType(raw string) OwnerType {
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
	ID        uuid.UUID
	Title     string
	OwnerID   string
	OwnerType OwnerType
	ChannelID string
	IsArchive bool
}

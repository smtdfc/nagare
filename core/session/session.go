package session

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/llm/provider"
)

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

type Info struct {
	ID        uuid.UUID
	Title     string
	OwnerID   string
	OwnerType OwnerType
	ChannelID string
	IsArchive bool

	CurrentLLMModel string
	LLMProviderID   uuid.UUID
	LLMProvider     *provider.LLMProviderInfo
}

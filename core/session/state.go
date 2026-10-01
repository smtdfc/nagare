package session

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type SessionState struct {
	SessionID          uuid.UUID
	CurrentModel       string
	CurrentLLMProvider uuid.UUID
	ChannelID          string
	Messages           messages.ListMessage
	NextCursor         string
	OwnerType          OwnerType
	OwnerID            string
}

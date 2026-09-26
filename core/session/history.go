package session

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/shared/messages"
)

type SessionHistory struct {
	SessionID  uuid.UUID
	ChannelID  string
	Messages   messages.ListMessage
	NextCursor string
	OwnerType  OwnerType
	OwnerID    string
}

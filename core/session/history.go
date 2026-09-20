package session

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/shared/message"
)

type SessionHistory struct {
	SessionID uuid.UUID
	ChannelID string
	Messages  message.ListMessage
	OwnerType OwnerType
	OwnerID   uuid.UUID
}

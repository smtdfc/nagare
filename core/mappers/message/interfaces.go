package message

import (
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type IMessageMapper interface {
	ToDomain(entity *entities.Message) (messages.Message, error)
	ToEntity(domain messages.Message, sessionID string) (*entities.Message, error)
	ToEntities(domains []messages.Message, sessionID string) ([]*entities.Message, error)
	ToDomains(entities []*entities.Message) ([]messages.Message, error)
}

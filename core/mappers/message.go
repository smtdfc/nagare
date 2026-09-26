package mappers

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/messages"
)

type MessageMapper struct {
}

func (m *MessageMapper) ToDomain(entity *entities.Message) (messages.Message, error) {
	if entity == nil {
		return nil, nil
	}

	switch entity.MessageKind {
	case string(messages.AgentStartedMessageType):
		return helpers.UnmarshalJson[messages.AgentStartedMessage](entity.Content)
	case string(messages.AgentCompletedMessageType):
		return helpers.UnmarshalJson[messages.AgentCompletedMessage](entity.Content)
	case string(messages.ReasoningMessageType):
		return helpers.UnmarshalJson[messages.ReasoningMessage](entity.Content)
	case string(messages.ResponseStartedMessageType):
		return helpers.UnmarshalJson[messages.ResponseStartedMessage](entity.Content)
	case string(messages.ResponseFailedMessageType):
		return helpers.UnmarshalJson[messages.ResponseFailedMessage](entity.Content)
	case string(messages.ResponseCompletedMessageType):
		return helpers.UnmarshalJson[messages.ResponseCompletedMessage](entity.Content)
	case string(messages.TextMessageType):
		return helpers.UnmarshalJson[messages.TextMessage](entity.Content)
	case string(messages.ToolCallMessageType):
		return helpers.UnmarshalJson[messages.ToolCallMessage](entity.Content)
	case string(messages.ToolResultMessageType):
		return helpers.UnmarshalJson[messages.ToolResultMessage](entity.Content)
	default:
		return nil, fmt.Errorf("failed to covert")
	}
}

func (m *MessageMapper) ToEntity(domain messages.Message, sessionID string) (*entities.Message, error) {
	if domain == nil {
		return nil, nil
	}

	id, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, err
	}

	raw, err := helpers.MarshalJson(domain)
	if err != nil {
		return nil, err
	}

	messageID, err := uuid.Parse(domain.GetMessageID())
	if err != nil {
		return nil, err
	}

	return &entities.Message{
		ID:          messageID,
		MessageKind: domain.GetMessageType().ToString(),
		Content:     raw,
		SessionID:   id,
	}, nil
}

func (m *MessageMapper) ToEntities(domains []messages.Message, sessionID string) ([]*entities.Message, error) {
	if domains == nil {
		return nil, nil
	}

	messageEntities := make([]*entities.Message, 0, len(domains))
	for _, d := range domains {
		if d == nil {
			continue
		}

		messageEntity, err := m.ToEntity(d, sessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to map domain to entity: %w", err)
		}

		messageEntities = append(messageEntities, messageEntity)
	}

	return messageEntities, nil
}

func (m *MessageMapper) ToDomains(entities []*entities.Message) ([]messages.Message, error) {
	if entities == nil {
		return nil, nil
	}

	domains := make([]messages.Message, 0, len(entities))
	for _, entity := range entities {
		if entity == nil {
			continue
		}

		domain, err := m.ToDomain(entity)
		if err != nil {
			return nil, fmt.Errorf("failed to map entity to domain for messages id %v: %w", entity.ID, err)
		}

		domains = append(domains, domain)
	}

	return domains, nil
}

// @Injectable
func NewMessageMapper() *MessageMapper {
	return &MessageMapper{}
}

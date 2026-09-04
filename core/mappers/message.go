package mappers

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/message"
)

type MessageMapper struct {
}

func (m *MessageMapper) ToDomain(entity *entities.Message) (message.Message, error) {
	if entity == nil {
		return nil, nil
	}

	switch entity.MessageKind {
	case string(message.AGENT_STARTED_MESSAGE):
		return helpers.UnmarshalJson[message.AgentStartedMessage](entity.Content)
	case string(message.AGENT_COMPLETED_MESSAGE):
		return helpers.UnmarshalJson[message.AgentCompletedMessage](entity.Content)
	case string(message.REASONING_MESSAGE):
		return helpers.UnmarshalJson[message.ReasoningMessage](entity.Content)
	case string(message.RESPONSE_STARTED_MESSAGE):
		return helpers.UnmarshalJson[message.ResponseStartedMessage](entity.Content)
	case string(message.RESPONSE_FAILED_MESSAGE):
		return helpers.UnmarshalJson[message.ResponseFailedMessage](entity.Content)
	case string(message.RESPONSE_COMPLETED_MESSAGE):
		return helpers.UnmarshalJson[message.ResponseCompletedMessage](entity.Content)
	case string(message.TEXT_MESSAGE):
		return helpers.UnmarshalJson[message.TextMessage](entity.Content)
	case string(message.TOOL_CALL_MESSAGE):
		return helpers.UnmarshalJson[message.ToolCallMessage](entity.Content)
	case string(message.TOOL_RESULT_MESSAGE):
		return helpers.UnmarshalJson[message.ToolResultMessage](entity.Content)
	default:
		return nil, fmt.Errorf("failed to covert")
	}
}

func (m *MessageMapper) ToEntity(domain message.Message, sessionID string) (*entities.Message, error) {
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

	return &entities.Message{
		MessageKind: domain.GetKind().ToString(),
		Content:     raw,
		SessionID:   id,
	}, nil
}

func (m *MessageMapper) ToDomains(entities []*entities.Message) ([]message.Message, error) {
	if entities == nil {
		return nil, nil
	}

	domains := make([]message.Message, 0, len(entities))
	for _, entity := range entities {
		if entity == nil {
			continue
		}

		domain, err := m.ToDomain(entity)
		if err != nil {
			return nil, fmt.Errorf("failed to map entity to domain for message id %v: %w", entity.ID, err)
		}

		domains = append(domains, domain)
	}

	return domains, nil
}

// @Injectable
func NewMessageMapper() *MessageMapper {
	return &MessageMapper{}
}

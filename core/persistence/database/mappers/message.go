package mappers

import (
	"fmt"

	"github.com/smtdfc/nagare/core/persistence/database/models"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/messages"
)

type MessageMapper struct{}

func (m *MessageMapper) ToDomain(model models.Message) (messages.Message, error) {
	switch model.MessageType {
	case string(messages.REASONING_MESSAGE):
		return helpers.MapObjectFromJson(model.Content, messages.NewReasoning(""))
	case string(messages.RESPONSE_STARTED_MESSAGE):
		return helpers.MapObjectFromJson(model.Content, messages.NewResponseStarted())
	case string(messages.RESPONSE_COMPLETED_MESSAGE):
		return helpers.MapObjectFromJson(model.Content, messages.NewResponseCompleted())
	case string(messages.RESPONSE_FAILED_MESSAGE):
		return helpers.MapObjectFromJson(model.Content, messages.NewResponseFailed("", ""))
	case string(messages.TEXT_MESSAGE):
		return helpers.MapObjectFromJson(model.Content, messages.NewText("", messages.SYSTEM))
	case string(messages.TOOL_CALL_MESSAGE):
		return helpers.MapObjectFromJson(model.Content, messages.NewToolCall("", "", ""))
	case string(messages.TOOL_RESULT_MESSAGE):
		return helpers.MapObjectFromJson(model.Content, messages.NewToolCallResult("", "", ""))
	default:
		return nil, fmt.Errorf("unknown message type: %s", model.MessageType)
	}
}

func (m *MessageMapper) ToModel(domain messages.Message) (*models.Message, error) {
	raw, err := helpers.MapObjectToJson(domain)
	if err != nil {
		return nil, err
	}
	var model = &models.Message{
		MessageType: string(domain.GetType()),
		Content:     raw,
	}

	return model, nil
}

package messages

import "github.com/smtdfc/nagare/shared/helpers"

type ReasoningMessage struct {
	ID      string      `json:"id"`
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
}

func (t *ReasoningMessage) GetMessageType() MessageType {
	return t.Type
}

func NewReasoningMessage(content string) *ReasoningMessage {
	return &ReasoningMessage{
		ID:      helpers.GenerateUUID(),
		Type:    ReasoningMessageType,
		Content: content,
	}
}

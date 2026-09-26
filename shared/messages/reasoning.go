package messages

import "github.com/smtdfc/nagare/shared/helpers"

type ReasoningMessage struct {
	ID      string      `json:"id"`
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
}

func (m *ReasoningMessage) GetMessageID() string {
	return m.ID
}

func (m *ReasoningMessage) GetMessageType() MessageType {
	return m.Type
}

func NewReasoningMessage(content string) *ReasoningMessage {
	return &ReasoningMessage{
		ID:      helpers.GenerateUUID(),
		Type:    ReasoningMessageType,
		Content: content,
	}
}

package messages

import "github.com/smtdfc/nagare/shared/helpers"

type TextMessage struct {
	ID      string      `json:"id"`
	Type    MessageType `json:"type"`
	Role    Role        `json:"role"`
	Content string      `json:"content"`
}

func (m *TextMessage) GetMessageID() string {
	return m.ID
}

func (m *TextMessage) GetMessageType() MessageType {
	return m.Type
}

func NewTextMessage(role Role, content string) *TextMessage {
	return &TextMessage{
		ID:      helpers.GenerateUUID(),
		Type:    TextMessageType,
		Role:    role,
		Content: content,
	}
}

package messages

import "github.com/smtdfc/nagare/shared/helpers"

type TextMessage struct {
	ID      string      `json:"id"`
	Type    MessageType `json:"type"`
	Role    Role        `json:"role"`
	Content string      `json:"content"`
}

func (t *TextMessage) GetMessageType() MessageType {
	return t.Type
}

func NewTextMessage(role Role, content string) *TextMessage {
	return &TextMessage{
		ID:      helpers.GenerateUUID(),
		Type:    TextMessageType,
		Role:    role,
		Content: content,
	}
}

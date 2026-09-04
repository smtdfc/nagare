package message

import "github.com/smtdfc/nagare/shared/helpers"

type TextMessage struct {
	ID      string `json:"id"`
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

func (t *TextMessage) GetKind() MessageKind {
	return TEXT_MESSAGE
}

func NewTextMessage(role Role, content string) *TextMessage {
	return &TextMessage{
		ID:      helpers.GenerateUUID(),
		Role:    role,
		Content: content,
	}
}

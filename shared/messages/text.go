package messages

import "github.com/google/uuid"

type Text struct {
	ID      string      `json:"id"`
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
	Role    Role        `json:"role"`
}

func NewText(content string, role Role) *Text {
	return &Text{
		ID:      uuid.New().String(),
		Type:    TEXT_MESSAGE,
		Content: content,
		Role:    role,
	}
}

func (t *Text) GetType() MessageType {
	return t.Type
}

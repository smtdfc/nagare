package messages

type Text struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
	Role    Role        `json:"role"`
}

func NewText(content string, role Role) *Text {
	return &Text{
		Type:    TEXT_MESSAGE,
		Content: content,
		Role:    role,
	}
}

func (t *Text) GetType() MessageType {
	return t.Type
}

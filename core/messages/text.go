package messages

type TextMessage struct {
	Role    Role
	Content string
}

func (m *TextMessage) GetType() string {
	return "Text"
}

package messages

type ReasoningMessage struct {
	Content string
}

func (m *ReasoningMessage) GetType() string {
	return "Reasoning"
}

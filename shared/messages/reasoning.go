package messages

type Reasoning struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
}

func (t *Reasoning) GetType() MessageType {
	return t.Type
}

func NewReasoning(t string) *Reasoning {
	return &Reasoning{
		Type:    REASONING_MESSAGE,
		Content: t,
	}
}

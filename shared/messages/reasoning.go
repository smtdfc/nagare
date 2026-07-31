package messages

import "github.com/google/uuid"

type Reasoning struct {
	ID      string      `json:"id"`
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
}

func (t *Reasoning) GetType() MessageType {
	return t.Type
}

func NewReasoning(t string) *Reasoning {
	return &Reasoning{
		ID:      uuid.New().String(),
		Type:    REASONING_MESSAGE,
		Content: t,
	}
}

package message

import "github.com/smtdfc/nagare/shared/helpers"

type ReasoningMessage struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func (t *ReasoningMessage) GetKind() MessageKind {
	return REASONING_MESSAGE
}

func NewReasoningMessage(content string) *ReasoningMessage {
	return &ReasoningMessage{
		ID:      helpers.GenerateUUID(),
		Content: content,
	}
}

package message

import "github.com/smtdfc/nagare/shared/helpers"

type ReasoningMessage struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func (t *ReasoningMessage) GetMessageType() MessageType {
	return ReasoningMessageType
}

func NewReasoningMessage(content string) *ReasoningMessage {
	return &ReasoningMessage{
		ID:      helpers.GenerateUUID(),
		Content: content,
	}
}

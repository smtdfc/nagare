package messages

import "github.com/smtdfc/nagare/pkgs/helpers"

type ReasoningMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invokeID"`
	Content  string      `json:"content"`
}

func (m *ReasoningMessage) GetMessageID() string {
	return m.ID
}

func (m *ReasoningMessage) GetMessageType() MessageType {
	return m.Type
}

func (m *ReasoningMessage) GetInvokeID() string {
	return m.InvokeID
}

func (m *ReasoningMessage) SetInvokeID(invokeID string) {
	m.InvokeID = invokeID
}

func NewReasoningMessage(content string) *ReasoningMessage {
	return &ReasoningMessage{
		ID:      helpers.GenerateUUID(),
		Type:    ReasoningMessageType,
		Content: content,
	}
}

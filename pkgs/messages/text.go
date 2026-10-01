package messages

import "github.com/smtdfc/nagare/pkgs/helpers"

type TextMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invoke_id"`
	Role     Role        `json:"role"`
	Content  string      `json:"content"`
}

func (m *TextMessage) GetMessageID() string {
	return m.ID
}

func (m *TextMessage) GetMessageType() MessageType {
	return m.Type
}

func (m *TextMessage) GetInvokeID() string {
	return m.InvokeID
}

func (m *TextMessage) SetInvokeID(invokeID string) {
	m.InvokeID = invokeID
}

func NewTextMessage(role Role, content string) *TextMessage {
	return &TextMessage{
		ID:      helpers.GenerateUUID(),
		Type:    TextMessageType,
		Role:    role,
		Content: content,
	}
}

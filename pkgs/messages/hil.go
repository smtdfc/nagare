package messages

import "github.com/smtdfc/nagare/pkgs/helpers"

type ApprovalRequestMessage struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	InvokeID  string      `json:"invokeID"`
	RequestID string      `json:"requestID"`
	SessionID string      `json:"sessionID"`
	Title     string      `json:"title"`
}

func (m *ApprovalRequestMessage) GetMessageID() string        { return m.ID }
func (m *ApprovalRequestMessage) GetMessageType() MessageType { return m.Type }
func (m *ApprovalRequestMessage) GetInvokeID() string         { return m.InvokeID }
func (m *ApprovalRequestMessage) SetInvokeID(invokeID string) { m.InvokeID = invokeID }

func NewApprovalRequestMessage(sessionID, requestID, title string) *ApprovalRequestMessage {
	return &ApprovalRequestMessage{
		ID:        helpers.GenerateUUID(),
		Type:      ApprovalRequestMessageType,
		RequestID: requestID,
		SessionID: sessionID,
		Title:     title,
	}
}

type ApprovalResponseMessage struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	InvokeID  string      `json:"invokeID"`
	RequestID string      `json:"requestID"`
	Accept    bool        `json:"accept"`
}

func (m *ApprovalResponseMessage) GetMessageID() string        { return m.ID }
func (m *ApprovalResponseMessage) GetMessageType() MessageType { return m.Type }
func (m *ApprovalResponseMessage) GetInvokeID() string         { return m.InvokeID }
func (m *ApprovalResponseMessage) SetInvokeID(invokeID string) { m.InvokeID = invokeID }

func NewApprovalResponseMessage(requestID string, accept bool) *ApprovalResponseMessage {
	return &ApprovalResponseMessage{
		ID:        helpers.GenerateUUID(),
		Type:      ApprovalResponseMessageType,
		RequestID: requestID,
		Accept:    accept,
	}
}

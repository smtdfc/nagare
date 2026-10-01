package messages

import "github.com/smtdfc/nagare/pkgs/helpers"

type ToolCallMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invoke_id"`
	CallID   string      `json:"call_id"`
	Name     string      `json:"name"`
	Args     string      `json:"args"`
}

func (m *ToolCallMessage) GetMessageID() string {
	return m.ID
}

func (m *ToolCallMessage) GetMessageType() MessageType {
	return m.Type
}

func (m *ToolCallMessage) GetInvokeID() string {
	return m.InvokeID
}

func (m *ToolCallMessage) SetInvokeID(invokeID string) {
	m.InvokeID = invokeID
}

func NewToolCallMessage(callID string, name string, args string) *ToolCallMessage {
	return &ToolCallMessage{
		ID:     helpers.GenerateUUID(),
		Type:   ToolCallMessageType,
		CallID: callID,
		Name:   name,
		Args:   args,
	}
}

type ToolResultMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invoke_id"`
	CallID   string      `json:"call_id"`
	Name     string      `json:"name"`
	Result   string      `json:"result"`
}

func (m *ToolResultMessage) GetMessageID() string {
	return m.ID
}

func (m *ToolResultMessage) GetMessageType() MessageType {
	return m.Type
}

func (m *ToolResultMessage) GetInvokeID() string {
	return m.InvokeID
}

func (m *ToolResultMessage) SetInvokeID(invokeID string) {
	m.InvokeID = invokeID
}

func NewToolResultMessage(callID string, name string, result string) *ToolResultMessage {
	return &ToolResultMessage{
		ID:     helpers.GenerateUUID(),
		Type:   ToolResultMessageType,
		CallID: callID,
		Name:   name,
		Result: result,
	}
}

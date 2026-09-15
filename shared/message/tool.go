package message

import "github.com/smtdfc/nagare/shared/helpers"

type ToolCallMessage struct {
	ID     string      `json:"id"`
	Type   MessageType `json:"type"`
	CallID string      `json:"call_id"`
	Name   string      `json:"name"`
	Args   string      `json:"args"`
}

func (t *ToolCallMessage) GetMessageType() MessageType {
	return t.Type
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
	ID     string      `json:"id"`
	Type   MessageType `json:"type"`
	CallID string      `json:"call_id"`
	Name   string      `json:"name"`
	Result string      `json:"result"`
}

func (t *ToolResultMessage) GetMessageType() MessageType {
	return t.Type
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

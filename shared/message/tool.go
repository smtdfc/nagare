package message

import "github.com/smtdfc/nagare/shared/helpers"

type ToolCallMessage struct {
	ID     string `json:"id"`
	CallID string `json:"call_id"`
	Name   string `json:"name"`
	Args   string `json:"args"`
}

func (t *ToolCallMessage) GetKind() Kind {
	return ToolCallMessageKind
}

func NewToolCallMessage(callID string, name string, args string) *ToolCallMessage {
	return &ToolCallMessage{
		ID:     helpers.GenerateUUID(),
		CallID: callID,
		Name:   name,
		Args:   args,
	}
}

type ToolResultMessage struct {
	ID     string `json:"id"`
	CallID string `json:"call_id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

func (t *ToolResultMessage) GetKind() Kind {
	return ToolResultMessageKind
}

func NewToolResultMessage(callID string, name string, result string) *ToolResultMessage {
	return &ToolResultMessage{
		ID:     helpers.GenerateUUID(),
		CallID: callID,
		Name:   name,
		Result: result,
	}
}

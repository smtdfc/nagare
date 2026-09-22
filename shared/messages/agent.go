package messages

import "github.com/smtdfc/nagare/shared/helpers"

type AgentStartedMessage struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
}

func (t *AgentStartedMessage) GetMessageType() MessageType {
	return t.Type
}

func NewAgentStartedMessage() *AgentStartedMessage {
	return &AgentStartedMessage{
		ID:   helpers.GenerateUUID(),
		Type: AgentStartedMessageType,
	}
}

type AgentCompletedMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	Success  bool        `json:"success"`
	Cancel   bool        `json:"cancel"`
	Duration float64     `json:"duration"`
}

func (t *AgentCompletedMessage) GetMessageType() MessageType {
	return t.Type
}

func NewAgentCompletedMessage(isSuccess bool, isCancel bool, duration float64) *AgentCompletedMessage {
	return &AgentCompletedMessage{
		ID:       helpers.GenerateUUID(),
		Type:     AgentCompletedMessageType,
		Success:  isSuccess,
		Cancel:   isCancel,
		Duration: duration,
	}
}

type AgentErrorMessage struct {
	ID    string      `json:"id"`
	Type  MessageType `json:"type"`
	Code  string      `json:"code"`
	Error string      `json:"error"`
}

func (t *AgentErrorMessage) GetMessageType() MessageType {
	return t.Type
}

func NewAgentErrorMessage(err string, code string) *AgentErrorMessage {
	return &AgentErrorMessage{
		ID:    helpers.GenerateUUID(),
		Type:  AgentErrorMessageType,
		Error: err,
		Code:  code,
	}
}

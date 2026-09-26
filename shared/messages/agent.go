package messages

import "github.com/smtdfc/nagare/shared/helpers"

type AgentStartedMessage struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
}

func (m *AgentStartedMessage) GetMessageID() string {
	return m.ID
}

func (m *AgentStartedMessage) GetMessageType() MessageType {
	return m.Type
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

func (m *AgentCompletedMessage) GetMessageID() string {
	return m.ID
}

func (m *AgentCompletedMessage) GetMessageType() MessageType {
	return m.Type
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

func (m *AgentErrorMessage) GetMessageID() string {
	return m.ID
}

func (m *AgentErrorMessage) GetMessageType() MessageType {
	return m.Type
}

func NewAgentErrorMessage(err string, code string) *AgentErrorMessage {
	return &AgentErrorMessage{
		ID:    helpers.GenerateUUID(),
		Type:  AgentErrorMessageType,
		Error: err,
		Code:  code,
	}
}

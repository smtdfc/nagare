package message

import "github.com/smtdfc/nagare/shared/helpers"

type AgentStartedMessage struct {
	ID string `json:"id"`
}

func (t *AgentStartedMessage) GetKind() MessageKind {
	return AGENT_STARTED_MESSAGE
}

func NewAgentStartedMessage() *AgentStartedMessage {
	return &AgentStartedMessage{
		ID: helpers.GenerateUUID(),
	}
}

type AgentCompletedMessage struct {
	ID       string  `json:"id"`
	Success  bool    `json:"success"`
	Cancel   bool    `json:"cancel"`
	Duration float64 `json:"duration"`
}

func (t *AgentCompletedMessage) GetKind() MessageKind {
	return AGENT_COMPLETED_MESSAGE
}

func NewAgentCompletedMessage(isSuccess bool, isCancel bool, duration float64) *AgentCompletedMessage {
	return &AgentCompletedMessage{
		ID:       helpers.GenerateUUID(),
		Success:  isSuccess,
		Cancel:   isCancel,
		Duration: duration,
	}
}

type AgentErrorMessage struct {
	ID    string `json:"id"`
	Code  string `json:"code"`
	Error string `json:"error"`
}

func (t *AgentErrorMessage) GetKind() MessageKind {
	return AGENT_ERROR_MESSAGE
}

func NewAgentErrorMessage(err string, code string) *AgentErrorMessage {
	return &AgentErrorMessage{
		ID:    helpers.GenerateUUID(),
		Error: err,
		Code:  code,
	}
}

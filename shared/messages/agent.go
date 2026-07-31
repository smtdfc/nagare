package messages

import "github.com/google/uuid"

type AgentResponseStatus string

const (
	AGENT_RESPONSE_COMPLETED AgentResponseStatus = "AGENT_RESPONSE_COMPLETED"
	AGENT_RESPONSE_STARTED   AgentResponseStatus = "AGENT_RESPONSE_STARTED"
	AGENT_RESPONSE_FAILED    AgentResponseStatus = "AGENT_RESPONSE_FAILED"
)

type AgentResponse struct {
	ID      string              `json:"id"`
	Type    MessageType         `json:"type"`
	Status  AgentResponseStatus `json:"status"`
	Content string              `json:"content"`
}

func (t *AgentResponse) GetType() MessageType {
	return t.Type
}

func NewAgentResponse(status AgentResponseStatus) *AgentResponse {
	return &AgentResponse{
		ID:     uuid.New().String(),
		Type:   AGENT_RESPONSE,
		Status: status,
	}
}

package messages

type AgentResponseDoneMessage struct{}

func (m *AgentResponseDoneMessage) GetType() string {
	return "AgentResponseDone"
}

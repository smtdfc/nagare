package agent

import "github.com/smtdfc/nagare/shared/messages"

type AgentState struct {
	Messages        messages.ListMessage
	PendingMessages messages.ListMessage
}

func (a *AgentState) WithHistory(history messages.ListMessage) *AgentState {
	a.Messages = append(a.Messages, history...)
	return a
}

func (a *AgentState) AddMessage(msg messages.Message) *AgentState {
	a.PendingMessages = append(a.PendingMessages, msg)
	return a
}

func (a *AgentState) GetHistory() messages.ListMessage {
	messages := make(messages.ListMessage, 0, len(a.Messages)+len(a.PendingMessages))
	messages = append(messages, a.Messages...)
	messages = append(messages, a.PendingMessages...)

	return messages
}

func (a *AgentState) CommitMessage() error {
	a.Messages = append(a.Messages, a.PendingMessages...)
	a.PendingMessages = messages.ListMessage{}
	return nil
}

func NewAgentState(listMessages messages.ListMessage) *AgentState {
	return &AgentState{
		Messages:        listMessages,
		PendingMessages: messages.ListMessage{},
	}
}

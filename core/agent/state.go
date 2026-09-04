package agent

import (
	"sync"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/shared/message"
)

type AgentState struct {
	mu             sync.RWMutex
	CurrentMessage message.ListMessage
	PendingMessage message.ListMessage
	ToolCalls      tool.ListToolCall
	LoopCounter    int
}

func (a *AgentState) GetFullMessage() message.ListMessage {
	a.mu.RLock()
	defer a.mu.RUnlock()

	messages := make(message.ListMessage, 0, len(a.CurrentMessage)+len(a.PendingMessage))
	messages = append(messages, a.CurrentMessage...)
	messages = append(messages, a.PendingMessage...)
	return messages
}

func (a *AgentState) SetMessages(messages message.ListMessage) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.CurrentMessage = messages
}

func (a *AgentState) AppendMessage(msg message.Message) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.PendingMessage = append(a.PendingMessage, msg)
}

func (a *AgentState) CommitMessage() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.CurrentMessage = append(a.CurrentMessage, a.PendingMessage...)
	a.PendingMessage = a.PendingMessage[:0]
}

func (a *AgentState) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.CurrentMessage = a.CurrentMessage[:0]
	a.PendingMessage = a.PendingMessage[:0]
	a.ToolCalls = a.ToolCalls[:0]
	a.LoopCounter = 0
}

func (a *AgentState) AddToolCall(toolCall *tool.ToolCall) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.ToolCalls = append(a.ToolCalls, toolCall)
}

func (a *AgentState) ResetToolCall() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.ToolCalls = a.ToolCalls[:0]
}

func (a *AgentState) IsToolCall() bool {
	return len(a.ToolCalls) > 0
}

func (a *AgentState) IncreaseLoopCounter() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.LoopCounter++
}

func (a *AgentState) GetLoopCounter() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.LoopCounter
}

func NewAgentState() *AgentState {
	return &AgentState{
		CurrentMessage: message.ListMessage{},
		PendingMessage: message.ListMessage{},
		ToolCalls:      tool.ListToolCall{},
		LoopCounter:    0,
	}
}

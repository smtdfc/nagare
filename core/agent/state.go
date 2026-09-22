package agent

import (
	"sync"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/shared/messages"
)

type State struct {
	mu             sync.RWMutex
	CurrentMessage messages.ListMessage
	PendingMessage messages.ListMessage
	ToolCalls      tool.ListToolCall
	LoopCounter    int
}

func (a *State) GetFullMessage() messages.ListMessage {
	a.mu.RLock()
	defer a.mu.RUnlock()

	messages := make(messages.ListMessage, 0, len(a.CurrentMessage)+len(a.PendingMessage))
	messages = append(messages, a.CurrentMessage...)
	messages = append(messages, a.PendingMessage...)
	return messages
}

func (a *State) SetMessages(messages messages.ListMessage) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.CurrentMessage = messages
}

func (a *State) AppendMessage(msg messages.Message) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.PendingMessage = append(a.PendingMessage, msg)
}

func (a *State) CommitMessage() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.CurrentMessage = append(a.CurrentMessage, a.PendingMessage...)
	a.PendingMessage = a.PendingMessage[:0]
}

func (a *State) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.CurrentMessage = a.CurrentMessage[:0]
	a.PendingMessage = a.PendingMessage[:0]
	a.ToolCalls = a.ToolCalls[:0]
	a.LoopCounter = 0
}

func (a *State) AddToolCall(toolCall *tool.ToolCall) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.ToolCalls = append(a.ToolCalls, toolCall)
}

func (a *State) ResetToolCall() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.ToolCalls = a.ToolCalls[:0]
}

func (a *State) IsToolCall() bool {
	return len(a.ToolCalls) > 0
}

func (a *State) IncreaseLoopCounter() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.LoopCounter++
}

func (a *State) GetLoopCounter() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.LoopCounter
}

func NewAgentState() *State {
	return &State{
		CurrentMessage: messages.ListMessage{},
		PendingMessage: messages.ListMessage{},
		ToolCalls:      tool.ListToolCall{},
		LoopCounter:    0,
	}
}

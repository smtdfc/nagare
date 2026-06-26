package agent

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/messages"
	"github.com/smtdfc/nagare/core/tool"
)

type AgentLoopState struct {
	Tools            domains.ListTool
	DynamicTools     domains.ListTool
	FinalTools       domains.ListTool
	History          messages.ListMessage
	IsDynamicToolSet bool
}

func (s *AgentLoopState) UseToolsForNextTurn(tools domains.ListTool) {
	s.DynamicTools = tools
	s.IsDynamicToolSet = true
}

func (s *AgentLoopState) InjectDynamicTool(t domains.Tool) {
	s.Tools = append(s.Tools, t)
}

func (s *AgentLoopState) BeforeTurn() {
	if s.IsDynamicToolSet {
		s.IsDynamicToolSet = false
		s.FinalTools = tool.MergeToolLists(s.Tools, s.DynamicTools)
	} else {
		s.FinalTools = s.Tools
	}
}

func (a *AgentLoopState) WithHistory(h messages.ListMessage) *AgentLoopState {
	a.History = make(messages.ListMessage, len(h))
	copy(a.History, h)
	return a
}

func (a *AgentLoopState) GetHistory(limit int) messages.ListMessage {
	start := 0
	if limit > 0 && len(a.History) > limit {
		start = len(a.History) - limit
	}

	relevantHistory := a.History[start:]
	result := make(messages.ListMessage, 0, len(relevantHistory)+2)

	result = append(result, SYSTEM_PROMPT)
	result = append(result, DEVELOPER_PROMPT)
	result = append(result, relevantHistory...)

	return result
}

func (a *AgentLoopState) ExtendHistory(history messages.ListMessage) *AgentLoopState {
	newHistory := make(messages.ListMessage, len(a.History)+len(history))
	copy(newHistory, a.History)
	copy(newHistory[len(a.History):], history)

	a.History = newHistory
	return a
}

func (a *AgentLoopState) AddHistory(msg messages.Message) {
	a.History = append(a.History, msg)
}

func (a *AgentLoopState) GetTools() domains.ListTool {
	return a.FinalTools
}

func NewAgentLoopState(toolReg *tool.ToolRegistry) *AgentLoopState {
	return &AgentLoopState{
		Tools:            toolReg.GetStaticTool(),
		DynamicTools:     domains.ListTool{},
		History:          DEFAULT_MESSAGES,
		IsDynamicToolSet: false,
	}
}

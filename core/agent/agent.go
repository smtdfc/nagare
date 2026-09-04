package agent

import (
	"context"

	"github.com/smtdfc/nagare/core/llm_provider"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/shared/message"
)

type Agent struct {
	toolMgr    *manager.ToolManager
	model      string
	llmAdapter llm_provider.LLMProviderAdapter
	executor   *AgentExecutor
	state      *AgentState
	logger     *logger.BaseLogger
}

func (a *Agent) Reset() *Agent {
	a.model = ""
	a.state.Reset()
	a.llmAdapter = nil
	return a
}

func (a *Agent) WithLLMAdapter(adapter llm_provider.LLMProviderAdapter) *Agent {
	a.llmAdapter = adapter
	return a
}

func (a *Agent) WithContext(messages message.ListMessage) *Agent {
	a.state.SetMessages(messages)
	return a
}

func (a *Agent) Invoke(ctx context.Context, msg message.Message, model string) (message.MessageChannel, error) {
	a.logger.Info("Start invoke agent")
	output := make(chan message.Message)
	go (func() {
		a.state.AppendMessage(msg)
		a.executor.Execute(ctx, model, a.llmAdapter, output)
	})()
	return output, nil
}

func NewAgent(toolMgr *manager.ToolManager, logger *logger.BaseLogger) *Agent {
	state := NewAgentState()
	return &Agent{
		model:      "",
		llmAdapter: nil,
		toolMgr:    toolMgr,
		executor:   NewAgentExecutor(state, toolMgr, logger),
		state:      state,
		logger:     logger.With("module", "agent"),
	}
}

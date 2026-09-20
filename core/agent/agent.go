package agent

import (
	"context"

	context2 "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/llm_provider"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/shared/message"
)

type InvokeOption struct {
	SessionID string
}

type Agent struct {
	toolMgr    *manager.ToolManager
	model      string
	llmAdapter llm_provider.LLMProviderAdapter
	executor   *Executor
	state      *State
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

func (a *Agent) Invoke(ctx context.Context, msg message.Message, model string, options *InvokeOption) (message.Channel, error) {
	a.logger.Info("Start invoke agent")
	output := make(chan message.Message)
	go (func() {
		a.state.AppendMessage(msg)
		ectx := &context2.ExecuteContext{
			Context: ctx,
		}
		if options != nil {
			ectx.SessionID = options.SessionID
		}
		
		a.executor.Execute(ectx, model, a.llmAdapter, output)
	})()
	return output, nil
}

func (a *Agent) DumpState() *State {
	return a.state
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

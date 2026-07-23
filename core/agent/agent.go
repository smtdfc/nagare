package agent

import (
	"strings"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/shared/messages"
)

type Agent struct {
	State       *AgentState
	LLMProvider domains.LLMProviderAdapter
	Model       string
	ToolMgr     *tool.ToolManager
}

func (a *Agent) WithLLMProvider(provider domains.LLMProviderAdapter) *Agent {
	a.LLMProvider = provider
	return a
}

func (a *Agent) WithModel(model string) *Agent {
	a.Model = model
	return a
}

func (a *Agent) WithState(state *AgentState) *Agent {
	a.State = state
	return a
}

func (a *Agent) WithToolManager(toolMgr *tool.ToolManager) *Agent {
	a.ToolMgr = toolMgr
	return a
}

func (a *Agent) Reset() {
	a.Model = ""
	a.LLMProvider = nil
	a.State = nil
}

func (a *Agent) Invoke(msg messages.Message) (domains.MessageChannel, error) {
	if a.LLMProvider == nil || a.State == nil {
		return nil, custom_errors.NewAgentError("Agent initialization failed. Please check the configuration settings")
	}

	ectx := context.NewExecuteContext(a.ToolMgr)
	a.State.AddMessage(msg)

	output := make(domains.MessageChannel)

	go func() {
		defer close(output)

		for {
			llmProviderOutput, _ := a.LLMProvider.Chat(a.Model, ectx, a.State.GetHistory())
			isFlushText := false
			var toolCalls = tool.ListToolCall{}

			var text strings.Builder
			var toolCallCount = 0
			for chunk := range llmProviderOutput {
				switch message := chunk.(type) {
				case *messages.Text:
					text.WriteString(message.Content)
					isFlushText = true
				case *messages.ToolCall:
					toolCallCount += 1
					a.State.AddMessage(&messages.ToolCall{
						Name:   message.Name,
						Args:   message.Args,
						CallID: message.CallID,
					})

					toolCalls = append(toolCalls, tool.NewToolCall(
						message.Name,
						message.Args,
						message.CallID,
					))

				default:
					if isFlushText {
						a.State.AddMessage(&messages.Text{
							Content: text.String(),
							Role:    messages.AGENT,
						})

						text.Reset()
					}
				}

				output <- chunk
			}

			if toolCallCount == 0 {
				break
			}

		}
	}()

	return output, nil
}

func NewAgent(model string, llmProvider domains.LLMProviderAdapter, state *AgentState) *Agent {
	return &Agent{
		Model:       model,
		State:       state,
		LLMProvider: llmProvider,
	}
}

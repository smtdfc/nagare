package agent

import (
	"fmt"
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
		output <- messages.NewAgentResponse(messages.AGENT_RESPONSE_STARTED)
		for {
			llmProviderOutput, err := a.LLMProvider.Chat(a.Model, ectx, a.State.GetHistory(), a.ToolMgr.GetListTool())
			if err != nil {
				msg := messages.NewAgentResponse(messages.AGENT_RESPONSE_FAILED)
				msg.Content = fmt.Sprintf("LLM Provider Error: %s", err.Error())
				output <- msg
				return
			}

			isFlushText := false
			var toolCalls = domains.ListToolCall{}

			var text strings.Builder
			var toolCallCount = 0
			for chunk := range llmProviderOutput {
				switch message := chunk.(type) {
				case *messages.Text:
					text.WriteString(message.Content)
					isFlushText = true
					output <- chunk
				case *messages.ToolCall:
					toolCallCount += 1
					a.State.AddMessage(messages.NewToolCall(message.Name, message.Args, message.CallID))
					toolCalls = append(toolCalls, domains.NewToolCall(
						message.Name,
						message.Args,
						message.CallID,
					))
					output <- chunk

				default:
					if isFlushText {
						a.State.AddMessage(messages.NewText(text.String(), messages.AGENT))
						text.Reset()
					}
				}
			}

			if toolCallCount == 0 {
				break
			}

			for _, call := range toolCalls {
				result := ectx.ExecuteToolCalls(call)
				if result.Status == domains.TOOL_CALL_PENDING {
					msg := messages.NewAgentResponse(messages.AGENT_RESPONSE_FAILED)
					msg.Content = "Execution race condition detected: Results were retrieved before the tool finished processing. Please verify that the tool operates synchronously."
					output <- msg
					return
				}

				a.State.AddMessage(messages.NewToolCallResult(
					result.CallID,
					result.Result,
					result.Error,
				))
			}
		}

		output <- messages.NewAgentResponse(messages.AGENT_RESPONSE_COMPLETED)
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

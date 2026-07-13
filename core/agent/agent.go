package agent

import (
	"strings"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/llm"
	"github.com/smtdfc/nagare/shared/messages"
)

type Agent struct {
	State       *AgentState
	LLMProvider llm.LLMProviderAdapter
	Model       string
}

func (a *Agent) Invoke(msg messages.Message) llm.MessageChannel {
	ectx := context.NewExecuteContext()
	a.State.AddMessage(msg)

	output := make(llm.MessageChannel)

	go func() {
		defer close(output)

		for {
			llmProviderOutput, _ := a.LLMProvider.Chat(a.Model, ectx, a.State.GetHistory())
			isFlushText := false
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

	return output
}
func NewAgent(model string, llmProvider llm.LLMProviderAdapter, state *AgentState) *Agent {
	return &Agent{
		Model:       model,
		State:       state,
		LLMProvider: llmProvider,
	}
}

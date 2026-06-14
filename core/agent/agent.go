package agent

import (
	"context"

	ectx "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/messages"
	"github.com/smtdfc/nagare/core/model"
	"github.com/smtdfc/nagare/core/tool"
)

type Agent struct {
	Model       model.ChatModel
	History     messages.ListMessage
	Tools       tool.ListTool
	Middlewares []any
}

func (a *Agent) WithHistory(h messages.ListMessage) *Agent {
	a.History = make(messages.ListMessage, len(h))
	copy(a.History, h)
	return a
}

func (a *Agent) ExtendHistory(history messages.ListMessage) *Agent {
	newHistory := make(messages.ListMessage, len(a.History)+len(history))
	copy(newHistory, a.History)
	copy(newHistory[len(a.History):], history)

	a.History = newHistory
	return a
}

func (a *Agent) WithTools(tools ...tool.Tool) *Agent {
	a.Tools = append(a.Tools, tools...)
	return a
}

func (a *Agent) Invoke(ctx context.Context, prompt string) <-chan messages.Message {
	ch := make(chan messages.Message)
	a.History = append(a.History, &messages.TextMessage{
		Role:    messages.USER,
		Content: prompt,
	})
	execCtx := ectx.NewExecuteContext(ctx, a.Tools)

	go func() {
		defer close(ch)
		a.processChat(execCtx, func(msg messages.Message) {
			ch <- msg
		})
	}()

	return ch
}

func (a *Agent) processChat(ctx ectx.ExecuteContext, cb model.MessageCallback) {
	for {
		fullTextMessage := ""
		var toolCalls []*messages.ToolCallMessage

		a.Model.Chat(ctx, a.History, func(msg messages.Message) {
			switch m := msg.(type) {
			case *messages.ToolCallMessage:
				toolCalls = append(toolCalls, m)
				a.History = append(a.History, m)
				cb(m)
			case *messages.TextMessage:
				cb(m)
				fullTextMessage += m.Content
			}
		})

		if fullTextMessage != "" {
			a.History = append(a.History, &messages.TextMessage{Role: messages.AGENT, Content: fullTextMessage})
		}

		if len(toolCalls) == 0 {
			break
		}

		for _, tc := range toolCalls {
			result, err := ctx.CallTool(tc.FunctionName, tc.Args)

			toolResultMsg := &messages.ToolResultMessage{
				CallID: tc.CallID,
				Result: result,
				Error:  err,
			}

			cb(toolResultMsg)
			a.History = append(a.History, toolResultMsg)
		}
	}
}

func NewAgent(model model.ChatModel) *Agent {
	return &Agent{
		Model:   model,
		History: messages.ListMessage{SYSTEM_PROMPT},
		Tools:   tool.ListTool{},
	}
}

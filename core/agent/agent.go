package agent

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ectx "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/messages"
	"github.com/smtdfc/nagare/core/model"
	"github.com/smtdfc/nagare/core/tool"
)

type Agent struct {
	Model       model.ChatModel
	History     messages.ListMessage
	Tools       tool.ListTool
	Middlewares []any
	logger      *slog.Logger
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
		err := a.processChat(execCtx, func(msg messages.Message) {
			ch <- msg
		})

		if err != nil {
			ch <- &messages.StreamErrorMessage{
				Cause: err.Error(),
			}
		}
	}()

	return ch
}

func (a *Agent) processChat(ctx ectx.ExecuteContext, cb model.MessageCallback) error {

	start := time.Now()
	a.logger.Info("Agent processing chat ")
	for {
		fullTextMessage := ""
		var toolCalls []*messages.ToolCallMessage

		err := a.Model.Chat(ctx, a.History, func(msg messages.Message) {
			switch m := msg.(type) {
			case *messages.ToolCallMessage:
				toolCalls = append(toolCalls, m)
				a.History = append(a.History, m)
				cb(m)
			case *messages.TextMessage:
				cb(m)
				fullTextMessage += m.Content
			case *messages.ReasoningMessage, *messages.ResponseFailedMessage, *messages.ResponseCreatedMessage, *messages.ResponseStatsMessage:
				cb(m)
			}
		})

		if err != nil {
			return err
		}

		if fullTextMessage != "" {
			a.History = append(a.History, &messages.TextMessage{Role: messages.AGENT, Content: fullTextMessage})
		}

		if len(toolCalls) == 0 {
			break
		}

		for _, tc := range toolCalls {
			a.logger.Info(fmt.Sprintf("Agent use tool %s", tc.FunctionName))
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

	cb(&messages.AgentResponseDoneMessage{})
	elapsed := time.Since(start)
	a.logger.Info(
		"Agent response completed",
		"duration", elapsed.String(),
	)
	return nil
}

func NewAgent(model model.ChatModel) *Agent {
	return &Agent{
		Model:   model,
		History: messages.ListMessage{SYSTEM_PROMPT},
		Tools:   tool.ListTool{},
		logger:  logger.GetLogger("Agent"),
	}
}

package agent

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ectx "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/messages"
	"github.com/smtdfc/nagare/core/model"
	"github.com/smtdfc/nagare/core/tool"
	nagare_logger "github.com/smtdfc/nagare/shared/logger"
)

type Agent struct {
	ToolRegistry *tool.ToolRegistry
	Model        model.ChatModel
	State        *AgentLoopState
	Middlewares  []any
	logger       *slog.Logger
}

func (a *Agent) Invoke(ctx context.Context, prompt string) <-chan messages.Message {
	ch := make(chan messages.Message)
	a.State.AddHistory(
		&messages.TextMessage{
			Role:    messages.USER,
			Content: prompt,
		},
	)

	execCtx := ectx.NewExecuteContext(ctx, a.ToolRegistry)

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
		a.State.BeforeTurn()
		fullTextMessage := ""
		var toolCalls []*messages.ToolCallMessage

		err := a.Model.Chat(ctx, a.State.GetHistory(NAGARE_LIST_MESSAGE_SIZE_LIMIT), func(msg messages.Message) {
			switch m := msg.(type) {
			case *messages.ToolCallMessage:
				toolCalls = append(toolCalls, m)
				a.State.History = append(a.State.History, m)
				cb(m)
			case *messages.TextMessage:
				cb(m)
				fullTextMessage += m.Content
			case *messages.ReasoningMessage, *messages.ResponseFailedMessage, *messages.ResponseCreatedMessage, *messages.ResponseStatsMessage:
				cb(m)
			}
		}, a.State.GetTools())

		if err != nil {
			return err
		}

		if fullTextMessage != "" {
			a.State.AddHistory(&messages.TextMessage{Role: messages.AGENT, Content: fullTextMessage})
		}

		if len(toolCalls) == 0 {
			break
		}

		for _, tc := range toolCalls {
			a.logger.Info(fmt.Sprintf("Agent use tool %s", tc.FunctionName))
			result, err := ctx.CallTool(tc.FunctionName, tc.Args)

			if result.Tool.GetType() == domains.DYNAMIC_TOOL && err == nil {
				a.State.InjectDynamicTool(result.Tool)
			}

			toolResultMsg := &messages.ToolResultMessage{
				CallID: tc.CallID,
				Result: result.Result,
				Error:  err,
			}

			cb(toolResultMsg)
			a.State.History = append(a.State.History, toolResultMsg)
		}

		ctx.CallMiddlewares(ectx.AFTER_TOOL_CALL, a.State)
	}

	cb(&messages.AgentResponseDoneMessage{})
	elapsed := time.Since(start)
	a.logger.Info(
		"Agent response completed",
		"duration", elapsed.String(),
	)
	return nil
}

func NewAgent(model model.ChatModel, toolReg *tool.ToolRegistry) *Agent {
	return &Agent{
		Model:        model,
		State:        NewAgentLoopState(toolReg),
		logger:       nagare_logger.GetLogger("Agent"),
		ToolRegistry: toolReg,
	}
}

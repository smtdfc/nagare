package agent

import (
	"strings"
	"time"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/llm_provider"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/message"
)

const MAX_AGENT_LOOP_COUNT = 20

type Executor struct {
	state   *State
	toolMgr *manager.ToolManager
	logger  *logger.BaseLogger
}

func (a *Executor) ExecuteTool(ctx *context.ExecuteContext, toolCall *tool.ToolCall) *tool.Result {
	result := a.toolMgr.Call(ctx, toolCall)
	if !result.IsSuccess {
		a.logger.Error("Failed to execute tool", "tool", toolCall.Name, "error", result.Result)
	} else {
		a.logger.Info("Execute tool", "tool", toolCall.Name)
	}

	return result
}

func (a *Executor) HandleError(message *message.ResponseFailedMessage) error {
	switch message.Code {
	case "429":
		return custom_errors.ErrModelQuotaExceed
	}

	return custom_errors.ErrUnknown
}

func (a *Executor) Execute(ctx *context.ExecuteContext, model string, llmAdapter llm_provider.LLMProviderAdapter, output message.WriteOnlyChannel) {
	defer close(output)

	isError := false
	isCancel := false
	var textBuilder strings.Builder
	isTextItem := false
	flushText := func() {
		if !isTextItem {
			return
		}

		a.state.AppendMessage(
			message.NewTextMessage(
				message.AGENT,
				textBuilder.String(),
			),
		)
		isTextItem = false
		textBuilder.Reset()
	}

	a.logger.Info("Start agent ReAct loop")
	startTime := time.Now()
	for {
		a.state.ResetToolCall()
		a.state.IncreaseLoopCounter()
		var executeError error
		select {
		case <-ctx.Done():
			a.logger.Error("Agent received cancel request")
			isCancel = true
		default:
		}

		if a.state.LoopCounter > MAX_AGENT_LOOP_COUNT {
			a.logger.Info("Agent loop exceeded maximum iterations")
			output <- message.NewAgentErrorMessage(
				custom_errors.ErrAgentExceedMaxIterations.Details,
				custom_errors.ErrAgentExceedMaxIterations.Code,
			)
			isCancel = true
		}

		if isCancel {
			break
		}

		llmOutput, err := llmAdapter.Send(ctx, model, a.state.GetFullMessage(), a.toolMgr.GetListTool())
		if err != nil {
			isError = true
			a.logger.Error("Agent error", "error", err)
			output <- message.NewAgentErrorMessage(err.Error(), custom_errors.ErrLLMProviderAdapter.Code)
			break
		}

		for chunk := range llmOutput {
			output <- chunk
			switch chunk.GetMessageType() {
			case message.TextMessageType:
				_, msg := helpers.SafeCast[*message.TextMessage](chunk)
				textBuilder.WriteString(msg.Content)
				isTextItem = true

			case message.ToolCallMessageType:
				flushText()
				_, msg := helpers.SafeCast[*message.ToolCallMessage](chunk)
				a.state.AddToolCall(tool.NewToolCall(
					msg.CallID,
					msg.Name,
					msg.Args,
				))
				a.state.AppendMessage(chunk)

			case message.ResponseFailedMessageType:
				flushText()
				_, msg := helpers.SafeCast[*message.ResponseFailedMessage](chunk)
				a.logger.Error("Agent error", "error", msg.Cause)
				executeError = a.HandleError(msg)
			}
		}

		flushText()

		if executeError != nil {
			isError = true
			a.logger.Info("Agent loop exited due to an error")
			output <- message.NewAgentErrorMessage(executeError.Error(), custom_errors.ErrAgentLoop.Code)
			break
		}

		if !a.state.IsToolCall() {
			break
		}

		for _, toolCall := range a.state.ToolCalls {
			result := a.ExecuteTool(ctx, toolCall)
			a.state.AppendMessage(result.ToMessage())
		}
	}

	endTime := time.Since(startTime).Seconds()
	output <- message.NewAgentCompletedMessage(
		!isError,
		isCancel,
		endTime,
	)
	a.logger.Info("Agent end loop", "duration", endTime)
}

func NewAgentExecutor(state *State, toolMgr *manager.ToolManager, logger *logger.BaseLogger) *Executor {
	return &Executor{
		state:   state,
		toolMgr: toolMgr,
		logger:  logger.With("module", "agent-executor"),
	}
}

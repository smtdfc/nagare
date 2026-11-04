package agent

import (
	"errors"
	"strings"
	"time"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/hil"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/message"
	"github.com/smtdfc/nagare/core/model_provider/provider"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/pkgs/helpers"
	"github.com/smtdfc/nagare/pkgs/messages"
)

const MaxAgentLoopCount = 20

type Executor struct {
	state       *State
	toolMgr     *manager.ToolManager
	approvalMgr *hil.Manager
	logger      *logger.BaseLogger
}

func (a *Executor) ExecuteTool(ctx *context.ExecuteContext, toolCall *tool.Call) *tool.Result {
	result := a.toolMgr.Call(ctx, toolCall)
	if !result.IsSuccess {
		a.logger.Error("Failed to execute tool", "tool", toolCall.Name, "error", result.Result)
	} else {
		a.logger.Info("Execute tool", "tool", toolCall.Name)
	}

	return result
}

func (a *Executor) HandleError(message *messages.ResponseFailedMessage) *custom_errors.NagareCoreError {
	cause := strings.TrimSpace(message.Cause)

	withCause := func(base *custom_errors.NagareCoreError) *custom_errors.NagareCoreError {
		if cause == "" {
			return base
		}

		return custom_errors.NewNagareCoreError(base.Code, cause)
	}

	switch message.Code {
	case "429":
		return withCause(custom_errors.ErrModelQuotaExceed)
	case "402":
		return withCause(custom_errors.ErrModelPaymentRequired)
	case "404":
		return withCause(custom_errors.ErrModelProviderNotFound)
	case "500", "502", "503", "504":
		return withCause(custom_errors.ErrModelProviderAdapter)
	}

	if cause != "" {
		return custom_errors.NewNagareCoreError(custom_errors.ErrUnknown.Code, cause)
	}

	return custom_errors.ErrUnknown
}

func (a *Executor) Execute(ctx *context.ExecuteContext, model string, modelAdapter provider.ModelProviderAdapter, output message.WriteOnlyChannel, opts *InvokeOption) {
	defer close(output)

	isError := false
	isCancel := false
	var textBuilder strings.Builder
	isTextItem := false
	flushText := func() {

		if !isTextItem {
			return
		}

		if !opts.DryMode {
			a.state.AppendMessage(
				messages.NewTextMessage(
					messages.AGENT,
					textBuilder.String(),
				),
			)
		}

		isTextItem = false
		textBuilder.Reset()
	}

	a.logger.Info("Start agent ReAct loop")
	startTime := time.Now()
	for {
		a.state.ResetToolCall()
		a.state.IncreaseLoopCounter()
		var executeError *custom_errors.NagareCoreError
		select {
		case <-ctx.Done():
			a.logger.Error("Agent received cancel request")
			isCancel = true
		default:
		}

		if a.state.LoopCounter > MaxAgentLoopCount {
			a.logger.Info("Agent loop exceeded maximum iterations")
			output <- messages.NewAgentErrorMessage(
				custom_errors.ErrAgentExceedMaxIterations.Details,
				custom_errors.ErrAgentExceedMaxIterations.Code,
			)
			isCancel = true
		}

		if isCancel {
			break
		}

		llmOutput, err := modelAdapter.Send(ctx, model, a.state.GetFullMessage(), a.toolMgr.GetListTool())
		if err != nil {
			if ctx.Err() != nil {
				isCancel = true
				break
			}
			isError = true
			a.logger.Error("Agent error", "error", err)

			code := custom_errors.ErrModelProviderAdapter.Code
			details := err.Error()
			var coreErr *custom_errors.NagareCoreError
			if errors.As(err, &coreErr) && coreErr != nil {
				code = coreErr.Code
				details = coreErr.Details
			}

			output <- messages.NewAgentErrorMessage(details, code)
			break
		}

		for chunk := range llmOutput {
			output <- chunk
			switch chunk.GetMessageType() {
			case messages.TextMessageType:
				_, msg := helpers.SafeCast[*messages.TextMessage](chunk)
				textBuilder.WriteString(msg.Content)
				isTextItem = true

			case messages.ToolCallMessageType:
				_, msg := helpers.SafeCast[*messages.ToolCallMessage](chunk)
				if opts.DryMode {
					continue
				}
				flushText()
				a.state.AddToolCall(tool.NewToolCall(
					msg.CallID,
					msg.Name,
					msg.Args,
				))
				a.state.AppendMessage(chunk)

			case messages.ResponseFailedMessageType:
				flushText()
				_, msg := helpers.SafeCast[*messages.ResponseFailedMessage](chunk)
				a.logger.Error("Agent error", "error", msg.Cause)
				executeError = a.HandleError(msg)
			}
		}

		flushText()

		if ctx.Err() != nil {
			isCancel = true
			break
		}

		if executeError != nil {
			isError = true
			a.logger.Info("Agent loop exited due to an error", "code", executeError.Code)
			output <- messages.NewAgentErrorMessage(executeError.Details, executeError.Code)
			break
		}

		if !a.state.IsToolCall() || opts.NoToolCall {
			break
		}

		for _, toolCall := range a.state.ToolCalls {
			if ctx.Err() != nil {
				isCancel = true
				break
			}
			result := a.ExecuteTool(ctx, toolCall)
			a.state.AppendMessage(result.ToMessage())
		}

		if isCancel {
			break
		}
	}

	endTime := time.Since(startTime).Seconds()
	output <- messages.NewAgentCompletedMessage(
		!isError,
		isCancel,
		endTime,
	)
	a.logger.Info("Agent end loop", "duration", endTime)
}

func NewAgentExecutor(state *State, toolMgr *manager.ToolManager, approvalMgr *hil.Manager, logger *logger.BaseLogger) *Executor {
	return &Executor{
		state:       state,
		toolMgr:     toolMgr,
		approvalMgr: approvalMgr,
		logger:      logger.With("module", "agent-executor"),
	}
}

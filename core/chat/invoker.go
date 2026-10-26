package chat

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/agent"
	config_mgr "github.com/smtdfc/nagare/core/config/manager"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/message"
	model_provider_mgr "github.com/smtdfc/nagare/core/model_provider/manager"
	"github.com/smtdfc/nagare/core/prompt"
	"github.com/smtdfc/nagare/core/session"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"
	tool_mgr "github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type AgentInvoker struct {
	agentPool        *agent.Pool
	sessionMgr       *session_mgr.SessionManager
	modelProviderMgr *model_provider_mgr.ModelProviderManager
	configMgr        *config_mgr.ConfigManager
	eventBus         *event_bus.CoreEventBus
	adapterLogger    *logger.BaseLogger
	toolMgr          *tool_mgr.ToolManager
	logger           *logger.BaseLogger
}

type AgentInvokeParams struct {
	SessionID        string
	InputMessages    messages.ListMessage
	SenderType       event_bus.SenderType
	SenderID         string
	SendIntoEventBus bool
	InvokeID         string
}

func (a *AgentInvoker) Invoke(ctx context.Context, params *AgentInvokeParams) (message.ReadOnlyChannel, error) {
	if params == nil {
		return nil, errors.New("agent invoke params cannot be nil")
	}

	output := make(chan messages.Message)
	sessionID := params.SessionID
	inputMessages := params.InputMessages
	senderType := params.SenderType
	senderID := params.SenderID
	sendIntoEventBus := params.SendIntoEventBus
	channelID := ""
	sessionOwnerType := ""
	sessionOwnerID := ""
	var sessionState *session.State
	invokeID := params.InvokeID
	if invokeID == "" {
		invokeID = uuid.New().String()
	}

	extractErrorDetails := func(err error) (string, string) {
		var coreErr *custom_errors.NagareCoreError
		if errors.As(err, &coreErr) && coreErr != nil {
			return coreErr.Code, coreErr.Details
		}
		return custom_errors.ErrAgentChat.Details, custom_errors.ErrAgentChat.Code
	}

	go func() {
		defer close(output)

		emit := func(msg messages.Message) {
			if msg != nil {
				msg.SetInvokeID(invokeID)
			}
			output <- msg
			if sendIntoEventBus && a.eventBus != nil {
				a.eventBus.Publish(ctx, event_bus.ChatMessageChunkEvent, &event_bus.ChatMessageChunkEventPayload{
					SessionID:        sessionID,
					ChannelID:        channelID,
					Chunk:            msg,
					SenderType:       senderType,
					SenderID:         senderID,
					SessionOwnerID:   sessionOwnerID,
					SessionOwnerType: sessionOwnerType,
				})
			}
		}

		systemPrompt, err := prompt.SystemPromptTemplate.Build(struct {
			ToolCategories string
		}{
			ToolCategories: a.toolMgr.GetCategoriesPrompt(),
		})
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}

		systemMessage := messages.NewTextMessage(messages.SYSTEM, systemPrompt)
		systemMessage.SetInvokeID(invokeID)
		history := messages.ListMessage{systemMessage}

		switch senderType {
		case event_bus.User:
			sessionState, err = a.sessionMgr.GetUserChatState(ctx, sessionID, senderID)
		case event_bus.Plugin:
			sessionState, err = a.sessionMgr.GetPluginChatState(ctx, sessionID, senderID)
		case event_bus.System:
			sessionState, err = a.sessionMgr.GetChatState(ctx, sessionID)
		}
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}
		channelID = sessionState.ChannelID
		history = append(history, sessionState.Messages...)
		sessionOwnerID = sessionState.OwnerID
		sessionOwnerType = sessionState.OwnerType.ToString()
		currentModelProviderID := sessionState.CurrentModelProvider
		currentModel := sessionState.CurrentModel

		if currentModelProviderID == uuid.Nil {
			emit(messages.NewAgentErrorMessage(
				custom_errors.ErrCurrentProviderNotSetup.Details,
				custom_errors.ErrCurrentProviderNotSetup.Code,
			))
			return
		}

		if currentModel == "" {
			emit(messages.NewAgentErrorMessage(
				custom_errors.ErrCurrentModelNotSetup.Details,
				custom_errors.ErrCurrentModelNotSetup.Code,
			))
			return
		}

		provider, err := a.modelProviderMgr.GetProviderByID(ctx, currentModelProviderID.String())
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}

		adapter, err := a.modelProviderMgr.GetAdapter(provider)
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}

		currentAgent := a.agentPool.Get().WithContext(history).WithModelAdapter(adapter)
		defer func() {
			currentAgent.Reset()
			a.agentPool.Put(currentAgent)
		}()

		for _, msg := range inputMessages {
			if msg != nil {
				msg.SetInvokeID(invokeID)
			}
		}
		agentOutput, err := currentAgent.Invoke(ctx, inputMessages, currentModel, &agent.InvokeOption{
			SessionID: sessionID,
		})
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}

		for msg := range agentOutput {
			emit(msg)
		}

		currentState := currentAgent.DumpState()
		for _, msg := range currentState.PendingMessage {
			if msg != nil {
				msg.SetInvokeID(invokeID)
			}
		}
		saveCtx, cancelSave := context.WithTimeout(context.WithoutCancel(ctx), 15*time.Second)
		defer cancelSave()
		err = a.sessionMgr.SaveHistory(saveCtx, sessionID, currentState.PendingMessage)
		if err != nil {
			a.logger.Error("failed to save history", "sessionID", sessionID, "err", err)
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}
	}()

	return output, nil
}

// @Injectable
func NewAgentInvoker(logger *logger.BaseLogger, agentPool *agent.Pool, modelProviderMgr *model_provider_mgr.ModelProviderManager, sessionMgr *session_mgr.SessionManager, configMgr *config_mgr.ConfigManager, eventBus *event_bus.CoreEventBus, toolMgr *tool_mgr.ToolManager) *AgentInvoker {
	return &AgentInvoker{
		agentPool:        agentPool,
		sessionMgr:       sessionMgr,
		configMgr:        configMgr,
		modelProviderMgr: modelProviderMgr,
		eventBus:         eventBus,
		logger:           logger.With("module", "agent-invoker"),
		adapterLogger:    logger.Clone(),
		toolMgr:          toolMgr,
	}
}

package chat

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/agent"
	config_mgr "github.com/smtdfc/nagare/core/config/manager"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/event_bus"
	llm_provider_mgr "github.com/smtdfc/nagare/core/llm_provider/manager"
	"github.com/smtdfc/nagare/core/logger"
	message "github.com/smtdfc/nagare/core/message"
	"github.com/smtdfc/nagare/core/prompt"
	"github.com/smtdfc/nagare/core/session"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"
	tool_mgr "github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/pkgs/messages"
)

type AgentInvoker struct {
	agentPool      *agent.Pool
	sessionMgr     *session_mgr.SessionManager
	llmProviderMgr *llm_provider_mgr.LLMProviderManager
	configMgr      *config_mgr.ConfigManager
	eventBus       *event_bus.CoreEventBus
	adapterLogger  *logger.BaseLogger
	toolMgr        *tool_mgr.ToolManager
	logger         *logger.BaseLogger
}

type AgentInvokeParams struct {
	SessionID        string
	InputMessages    messages.ListMessage
	SenderType       event_bus.SenderType
	SenderID         string
	SendIntoEventBus bool
	InvokeID         string
}

func (a *AgentInvoker) Invoke(params *AgentInvokeParams) (message.ReadOnlyChannel, error) {
	if params == nil {
		return nil, errors.New("agent invoke params cannot be nil")
	}

	output := make(chan messages.Message)
	ctx := context.Background()
	sessionID := params.SessionID
	inputMessages := params.InputMessages
	senderType := params.SenderType
	senderID := params.SenderID
	sendIntoEventBus := params.SendIntoEventBus
	channelID := ""
	sessionOwnerType := ""
	sessionOwnerID := ""
	var sessionState *session.SessionState
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
				a.eventBus.Publish(ctx, event_bus.ChunkEvent, &event_bus.ChatChunkEventPayload{
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
		currentLLMProviderID := sessionState.CurrentLLMProvider
		currentLLMModel := sessionState.CurrentModel

		if currentLLMProviderID == uuid.Nil {
			emit(messages.NewAgentErrorMessage(
				custom_errors.ErrCurrentProviderNotSetup.Details,
				custom_errors.ErrCurrentProviderNotSetup.Code,
			))
			return
		}

		if currentLLMModel == "" {
			emit(messages.NewAgentErrorMessage(
				custom_errors.ErrCurrentModelNotSetup.Details,
				custom_errors.ErrCurrentModelNotSetup.Code,
			))
			return
		}

		provider, err := a.llmProviderMgr.GetProviderByID(ctx, currentLLMProviderID.String())
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}

		adapter, err := a.llmProviderMgr.GetAdapter(provider)
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}

		currentAgent := a.agentPool.Get().WithContext(history).WithLLMAdapter(adapter)
		for _, msg := range inputMessages {
			if msg != nil {
				msg.SetInvokeID(invokeID)
			}
		}
		agentOutput, err := currentAgent.Invoke(ctx, inputMessages, currentLLMModel, &agent.InvokeOption{
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

		defer func() {
			currentAgent.Reset()
			a.agentPool.Put(currentAgent)
		}()

		currentState := currentAgent.DumpState()
		for _, msg := range currentState.PendingMessage {
			if msg != nil {
				msg.SetInvokeID(invokeID)
			}
		}
		err = a.sessionMgr.SaveHistory(ctx, sessionID, currentState.PendingMessage)
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}
	}()

	return output, nil
}

// @Injectable
func NewAgentInvoker(logger *logger.BaseLogger, agentPool *agent.Pool, llmProviderMgr *llm_provider_mgr.LLMProviderManager, sessionMgr *session_mgr.SessionManager, configMgr *config_mgr.ConfigManager, eventBus *event_bus.CoreEventBus, toolMgr *tool_mgr.ToolManager) *AgentInvoker {
	return &AgentInvoker{
		agentPool:      agentPool,
		sessionMgr:     sessionMgr,
		configMgr:      configMgr,
		llmProviderMgr: llmProviderMgr,
		eventBus:       eventBus,
		logger:         logger.With("module", "agent-invoker"),
		adapterLogger:  logger.Clone(),
		toolMgr:        toolMgr,
	}
}

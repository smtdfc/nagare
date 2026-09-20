package chat

import (
	"context"
	"errors"

	"github.com/smtdfc/nagare/core/agent"
	config_mgr "github.com/smtdfc/nagare/core/config/manager"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/event_bus"
	llm_provider_mgr "github.com/smtdfc/nagare/core/llm_provider/manager"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/session"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"

	"github.com/smtdfc/nagare/shared/message"
)

type AgentInvoker struct {
	agentPool      *agent.Pool
	sessionMgr     *session_mgr.SessionManager
	llmProviderMgr *llm_provider_mgr.LLMProviderManager
	configMgr      *config_mgr.ConfigManager
	eventBus       *event_bus.CoreEventBus
	adapterLogger  *logger.BaseLogger
	logger         *logger.BaseLogger
}

func (a *AgentInvoker) Invoke(
	sessionID string,
	text string,
	senderType event_bus.SenderType,
	senderId string,
	sendIntoEventBus bool,
) (message.ReadOnlyChannel, error) {
	output := make(chan message.Message)
	ctx := context.Background()
	channelID := ""
	sessionOwnerType := ""
	sessionOwnerID := ""
	var sessionHistory *session.SessionHistory
	var err error

	extractErrorDetails := func(err error) (string, string) {
		var coreErr *custom_errors.NagareCoreError
		if errors.As(err, &coreErr) && coreErr != nil {
			return coreErr.Code, coreErr.Details
		}
		return custom_errors.ErrAgentChat.Details, custom_errors.ErrAgentChat.Code
	}

	go func() {
		defer close(output)

		emit := func(msg message.Message) {
			output <- msg
			if sendIntoEventBus && a.eventBus != nil {
				a.eventBus.Publish(ctx, event_bus.ChunkEvent, &event_bus.ChatChunkEventPayload{
					SessionID:        sessionID,
					ChannelID:        channelID,
					Chunk:            msg,
					SenderType:       senderType,
					SenderID:         senderId,
					SessionOwnerID:   sessionOwnerID,
					SessionOwnerType: sessionOwnerType,
				})
			}
		}

		history := make([]message.Message, 0)
		history = append(history, message.NewTextMessage(
			message.DEVELOPER, `
			- You MUST reply using the EXACT same language that the user is currently using in their prompt/request.
			- DO NOT treat this message as input, a question, or a command from the User.
			- DO NOT attempt to create new tasks, DO NOT ask the user for more details/due dates, and DO NOT schedule anything.
			- ABSOLUTELY FORBIDDEN to reply with generic assistant fluff like "Got it! I’ll set up a reminder for you...", "Understood, I will remind you...", or any similar nonsense.
			- DO NOT output raw tool calls, function execution JSON, or technical diagnostic data to the user. Process them internally and reply only with the final natural language response.
		`))

		if senderType == event_bus.User {
			sessionHistory, err = a.sessionMgr.GetUserChatHistory(ctx, sessionID)
		} else if senderType == event_bus.Plugin {
			sessionHistory, err = a.sessionMgr.GetPluginChatHistory(ctx, sessionID, senderId)
		} else if senderType == event_bus.System {
			sessionHistory, err = a.sessionMgr.GetChatHistory(ctx, sessionID)
		}
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(message.NewAgentErrorMessage(details, code))
			return
		}
		channelID = sessionHistory.ChannelID
		history = append(history, sessionHistory.Messages...)
		sessionOwnerID = sessionHistory.OwnerID.String()
		sessionOwnerType = sessionHistory.OwnerType.ToString()

		generalConfig, err := a.configMgr.GetGeneralConfig(ctx)
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(message.NewAgentErrorMessage(details, code))
			return
		}

		if generalConfig.CurrentProvider == "" {
			emit(message.NewAgentErrorMessage(
				custom_errors.ErrCurrentProviderNotSetup.Details,
				custom_errors.ErrCurrentProviderNotSetup.Code,
			))
			return
		}

		if generalConfig.CurrentModel == "" {
			emit(message.NewAgentErrorMessage(
				custom_errors.ErrCurrentModelNotSetup.Details,
				custom_errors.ErrCurrentModelNotSetup.Code,
			))
			return
		}

		provider, err := a.llmProviderMgr.GetProviderByID(ctx, generalConfig.CurrentProvider)
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(message.NewAgentErrorMessage(details, code))
			return
		}

		adapter, err := a.llmProviderMgr.GetAdapter(provider)
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(message.NewAgentErrorMessage(details, code))
			return
		}

		currentAgent := a.agentPool.Get().WithContext(history).WithLLMAdapter(adapter)
		agentOutput, err := currentAgent.Invoke(ctx, message.NewTextMessage(
			message.USER,
			text,
		), generalConfig.CurrentModel, &agent.InvokeOption{
			SessionID: sessionID,
		})
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(message.NewAgentErrorMessage(details, code))
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
		err = a.sessionMgr.SaveHistory(ctx, sessionID, currentState.PendingMessage)
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(message.NewAgentErrorMessage(details, code))
			return
		}
	}()

	return output, nil
}

// @Injectable
func NewAgentInvoker(logger *logger.BaseLogger, agentPool *agent.Pool, llmProviderMgr *llm_provider_mgr.LLMProviderManager, sessionMgr *session_mgr.SessionManager, configMgr *config_mgr.ConfigManager, eventBus *event_bus.CoreEventBus) *AgentInvoker {
	return &AgentInvoker{
		agentPool:      agentPool,
		sessionMgr:     sessionMgr,
		configMgr:      configMgr,
		llmProviderMgr: llmProviderMgr,
		eventBus:       eventBus,
		logger:         logger.With("module", "agent-invoker"),
		adapterLogger:  logger.Clone(),
	}
}

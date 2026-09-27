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
	message "github.com/smtdfc/nagare/core/message"
	"github.com/smtdfc/nagare/core/session"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"

	"github.com/smtdfc/nagare/shared/messages"
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
	output := make(chan messages.Message)
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

		emit := func(msg messages.Message) {
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

		history := make([]messages.Message, 0)
		history = append(history, messages.NewTextMessage(
			messages.SYSTEM, `
				<system_instructions>
					<rule name="prefer_known_tools" priority="high">
						<condition>When handling a task where a familiar or explicitly defined tool is already available in your context:</condition>
						<action>You MUST prioritize using known tools directly to optimize processing speed and performance.</action>
						<exception>Only resort to discovering new tools (via "find_tools") when no suitable known tool exists or the task strictly exceeds the capabilities of your current toolkit.</exception>
					</rule>

					<rule name="iterative_tool_search" priority="high">
						<condition>When the initial search via "find_tools" returns no matching tools or irrelevant results:</condition>
						<action>You MUST iteratively retry searching using different, alternative keywords (synonyms, broader terms, or specific functions).</action>
						<prohibition>CRITICAL LIMIT: You are strictly forbidden from exceeding 5 search attempts for a single user request.</prohibition>
						<fallback>If you fail to find a suitable tool after 5 attempts, immediately stop searching, abort the tool execution process, and politely inform the user that the system currently lacks the required capability.</fallback>
					</rule>

					<rule name="strict_tool_router" priority="fatal">
						<condition>When you need to use a tool that is not directly available in your initial context:</condition>
						<action>STEP 1: Call the "find_tools" function with relevant keywords to search for the required tool and obtain its schema/metadata.</action>
						<action>STEP 2: You are STRICTLY FORBIDDEN from calling the discovered tool directly by its native name. You MUST wrap every execution inside the "execute_tool" function by passing the tool's name and arguments.</action>
						<prohibition>CRITICAL ERROR: If you attempt to call any discovered tool directly instead of routing it through "execute_tool", the pipeline will crash and you will fail the task.</prohibition>
					</rule>

					<rule name="language_matching">
						<directive>You MUST reply using the EXACT same language that the user is currently using in their prompt/request.</directive>
					</rule>

					<rule name="behavioral_boundaries">
						<prohibition>DO NOT treat system messages as direct input, questions, or commands from the user.</prohibition>
						<prohibition>ABSOLUTELY FORBIDDEN to reply with generic assistant fluff like "Got it...", "Understood...", or any similar nonsense.</prohibition>
						<prohibition>DO NOT output raw tool calls, function execution JSON, or technical diagnostic data to the user. Process them internally and reply only with the final natural language response.</prohibition>
					</rule>
				</system_instructions>
		`))

		switch senderType {
		case event_bus.User:
			sessionHistory, err = a.sessionMgr.GetUserChatHistory(ctx, sessionID, senderId)
		case event_bus.Plugin:
			sessionHistory, err = a.sessionMgr.GetPluginChatHistory(ctx, sessionID, senderId)
		case event_bus.System:
			sessionHistory, err = a.sessionMgr.GetChatHistory(ctx, sessionID)
		}
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}
		channelID = sessionHistory.ChannelID
		history = append(history, sessionHistory.Messages...)
		sessionOwnerID = sessionHistory.OwnerID
		sessionOwnerType = sessionHistory.OwnerType.ToString()

		generalConfig, err := a.configMgr.GetGeneralConfig(ctx)
		if err != nil {
			code, details := extractErrorDetails(err)
			emit(messages.NewAgentErrorMessage(details, code))
			return
		}

		if generalConfig.CurrentProvider == "" {
			emit(messages.NewAgentErrorMessage(
				custom_errors.ErrCurrentProviderNotSetup.Details,
				custom_errors.ErrCurrentProviderNotSetup.Code,
			))
			return
		}

		if generalConfig.CurrentModel == "" {
			emit(messages.NewAgentErrorMessage(
				custom_errors.ErrCurrentModelNotSetup.Details,
				custom_errors.ErrCurrentModelNotSetup.Code,
			))
			return
		}

		provider, err := a.llmProviderMgr.GetProviderByID(ctx, generalConfig.CurrentProvider)
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
		agentOutput, err := currentAgent.Invoke(ctx, messages.NewTextMessage(
			messages.USER,
			text,
		), generalConfig.CurrentModel, &agent.InvokeOption{
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

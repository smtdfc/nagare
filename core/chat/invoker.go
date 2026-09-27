package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/smtdfc/nagare/core/agent"
	config_mgr "github.com/smtdfc/nagare/core/config/manager"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/event_bus"
	llm_provider_mgr "github.com/smtdfc/nagare/core/llm_provider/manager"
	"github.com/smtdfc/nagare/core/logger"
	message "github.com/smtdfc/nagare/core/message"
	"github.com/smtdfc/nagare/core/session"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"
	tool_mgr "github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/shared/messages"
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

		toolCategories := a.toolMgr.GetCategoriesString()
		history := make([]messages.Message, 0)
		history = append(history, messages.NewTextMessage(
			messages.DEVELOPER,
			fmt.Sprintf(`
				<system_instructions>
					<metadata>
						<attr key="name">Nagare</attr>
						<attr key="description">Virtual assistant operating on the computer</attr>
					</metadata>
					<rule name="tool_routing">
						<description>Automated process for discovering, selecting, and executing tools when a request lacks the necessary handling capability.</description>
						<conditions>
							<condition>Activate only when the user request exceeds current capabilities and requires external tool assistance.</condition>
							<condition>Input categories (%s) must be thoroughly analyzed to precisely match the functionality of the tools.</condition>
						</conditions>

						<steps>
							<step n="1">
								<action>Analyze syntax and semantics of the user request.</action>
								<details>Determine the exact category or domain from the permitted set: %s.</details>
							</step>

							<step n="2">
								<action>Invoke the system discovery tool.</action>
								<details>Execute the "find_tool_by_categories" function with the determined category from Step 1 as a parameter to retrieve a list of available tools.</details>
							</step>

							<step n="3">
								<action>Evaluate and select the optimal tool.</action>
								<details>Carefully review the returned list, compare the features of each tool, and select the most appropriate one matching the user's intent. If no matching tool is found, halt the process and notify the user.</details>
							</step>

							<step n="4">
								<action>Configure parameters and execute.</action>
								<details>Prepare all required parameters according to the selected tool's schema, then call the "execute_tool" function to run it. Ensure no parameter is missing or of the wrong data type to prevent system errors.</details>
							</step>
						</steps>
						
						<fallback>
							If "execute_tool" returns an error due to invalid parameters, the system must automatically review the parameter structure in Step 4, correct the error, and retry execution at most once before reporting an error to the user.
						</fallback>
					</rule>
					<rule name="response_language">
						<description>Rule for controlling and maintaining the assistant's response language.</description>
						<conditions>
							<condition>Must strictly adhere to the language currently being used by the user in the conversation.</condition>
						</conditions>
						<steps>
							<step n="1">
								<action>Identify the user's language.</action>
								<details>Analyze the latest input message to accurately recognize the language or terminology used by the user.</details>
							</step>
							<step n="2">
								<action>Format the output language.</action>
								<details>The entire response content must be written completely in the language identified in Step 1.</details>
							</step>
						</steps>
						<constraints>
							<constraint>Do not arbitrarily switch to another language (e.g., automatically switching from Vietnamese to English or vice versa) unless explicitly requested by the user.</constraint>
						</constraints>
					</rule>
				</system_instructions>
			`, toolCategories, toolCategories,
			),
		))

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

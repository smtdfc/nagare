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

		toolCategories := a.toolMgr.GetCategoriesPrompt()
		fmt.Println(toolCategories)
		history := make([]messages.Message, 0)
		history = append(history, messages.NewTextMessage(
			messages.SYSTEM,
			fmt.Sprintf(`
				<system_instructions>
					<metadata>
						<attr key="name">Nagare</attr>
						<attr key="description">Virtual assistant operating on the computer</attr>
					</metadata>
					<rule name="tool_routing">
						<description>
							Deterministic procedure for discovering, selecting, and executing
							tools required to fulfill a user request.
						</description>

						<constraints>
							<constraint>
								Any request that requires an external tool MUST enter this routing process.
							</constraint>

							<constraint>
								The assistant MUST NOT claim that a suitable tool is unavailable
								before executing "find_tool_by_categories".
							</constraint>

							<constraint>
								The assistant MUST NOT skip the discovery step for a tool-dependent request.
							</constraint>
						</constraints>
						<categories>
							%s
						</categories>
						<steps>
							<step id="analyze_request">
								<action>Determine relevant tool categories.</action>
								<details>
									Analyze the user's intent and select one or more categories from the available set.
								</details>
							</step>

							<step id="discover_tools">
								<action>Discover available tools.</action>
								<details>
									MUST execute "find_tool_by_categories" using the categories
									determined in the previous step.
								</details>
							</step>

							<step id="select_tool">
								<action>Select the appropriate tool.</action>
								<details>
									Review the tools returned by "find_tool_by_categories"
									and select the tool that best matches the user's intent.

									If multiple tools are required, determine the appropriate
									execution order.

									If no suitable tool is returned, report that the requested
									capability is unavailable.
								</details>
							</step>

							<step id="execute_tool">
								<action>Configure and execute the selected tool.</action>
								<details>
									Construct the required parameters according to the selected
									tool's schema, then execute "execute_tool".

									Ensure that all required parameters are present and that
									their values match the expected data types.
								</details>
							</step>
						</steps>

						<fallback>
							<condition>
								"execute_tool" returns an error caused by invalid parameters.
							</condition>

							<action>
								Re-check the selected tool's schema, correct the invalid
								parameters, and retry execution once.
							</action>

							<failure>
								If the retry fails, report the execution error to the user.
							</failure>
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
			`, toolCategories,
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

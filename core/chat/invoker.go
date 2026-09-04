package chat

import (
	"context"
	"errors"

	"github.com/smtdfc/nagare/core/agent"
	config_mgr "github.com/smtdfc/nagare/core/config/manager"
	"github.com/smtdfc/nagare/core/custom_errors"
	llm_provider_mgr "github.com/smtdfc/nagare/core/llm_provider/manager"
	"github.com/smtdfc/nagare/core/logger"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"

	"github.com/smtdfc/nagare/shared/message"
)

type AgentInvoker struct {
	agentPool      *agent.AgentPool
	sessionMgr     *session_mgr.SessionManager
	llmProviderMgr *llm_provider_mgr.LLMProviderManager
	configMgr      *config_mgr.ConfigManager
	adapterLogger  *logger.BaseLogger
	logger         *logger.BaseLogger
}

func (a *AgentInvoker) Invoke(
	sessionID string,
	text string,
) (message.MessageReadOnlyChannel, error) {
	output := make(chan message.Message)
	ctx := context.Background()

	extractErrorDetails := func(err error) (string, string) {
		var coreErr *custom_errors.NagareCoreError
		if errors.As(err, &coreErr) && coreErr != nil {
			return coreErr.Code, coreErr.Details
		}
		return custom_errors.ErrAgentChat.Details, custom_errors.ErrAgentChat.Code
	}

	go func() {
		defer close(output)

		generalConfig, err := a.configMgr.GetGeneralConfig(ctx)
		if err != nil {
			code, details := extractErrorDetails(err)
			output <- message.NewAgentErrorMessage(details, code)
			return
		}

		history, err := a.sessionMgr.GetUserChatHistory(ctx, sessionID)
		if err != nil {
			code, details := extractErrorDetails(err)
			output <- message.NewAgentErrorMessage(details, code)
			return
		}

		if generalConfig.CurrentProvider == "" {
			output <- message.NewAgentErrorMessage(
				custom_errors.ErrCurrentProviderNotSetup.Details,
				custom_errors.ErrCurrentProviderNotSetup.Code,
			)
			return
		}

		if generalConfig.CurrentModel == "" {
			output <- message.NewAgentErrorMessage(
				custom_errors.ErrCurrentModelNotSetup.Details,
				custom_errors.ErrCurrentModelNotSetup.Code,
			)
			return
		}

		provider, err := a.llmProviderMgr.GetProviderByID(ctx, generalConfig.CurrentProvider)
		if err != nil {
			code, details := extractErrorDetails(err)
			output <- message.NewAgentErrorMessage(details, code)
			return
		}

		adapter, err := a.llmProviderMgr.GetAdapter(provider)
		if err != nil {
			code, details := extractErrorDetails(err)
			output <- message.NewAgentErrorMessage(details, code)
			return
		}

		agent := a.agentPool.Get().WithContext(history).WithLLMAdapter(adapter)
		agentOutput, err := agent.Invoke(ctx, message.NewTextMessage(
			message.USER,
			text,
		), generalConfig.CurrentModel)
		if err != nil {
			code, details := extractErrorDetails(err)
			output <- message.NewAgentErrorMessage(details, code)
			return
		}
		for msg := range agentOutput {
			output <- msg
		}

		agent.Reset()
		a.agentPool.Put(agent)
	}()

	return output, nil
}

// @Injectable
func NewAgentInvoker(logger *logger.BaseLogger, agentPool *agent.AgentPool, llmProviderMgr *llm_provider_mgr.LLMProviderManager, sessionMgr *session_mgr.SessionManager, configMgr *config_mgr.ConfigManager) *AgentInvoker {
	return &AgentInvoker{
		agentPool:      agentPool,
		sessionMgr:     sessionMgr,
		configMgr:      configMgr,
		llmProviderMgr: llmProviderMgr,
		logger:         logger.With("module", "agent-invoker"),
		adapterLogger:  logger.Clone(),
	}
}

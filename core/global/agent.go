package global

import (
	"errors"

	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/llm/providers"
	"github.com/smtdfc/nagare/shared/messages"
)

const AGENT_POOL_SIZE = 10

var GlobalAgentPool *agent.AgentPool

func CreateEmptyAgentState() *agent.AgentState {
	return agent.NewAgentState(messages.EMPTY_LIST)
}

func GetAgentFromPool() *agent.Agent {
	return GlobalAgentPool.Get()
}

func PutAgentIntoPool(agent *agent.Agent) {
	GlobalAgentPool.Put(agent)
}

func FetchReadyAgent(state *agent.AgentState) (*agent.Agent, error) {

	generalConfig, err := GlobalConfigMgr.GetGeneralConfig()
	if err != nil {
		return nil, err
	}

	if generalConfig.CurrentProvider == "" {
		return nil, errors.New("LLM provider not configured. Please configure it first.")
	}

	currentProviderConfig, err := GlobalConfigMgr.GetLLMProviderConfigByID(generalConfig.CurrentProvider)
	if err != nil {
		return nil, err
	}

	if currentProviderConfig == nil {
		return nil, errors.New("LLM provider not found. Please configure it first.")
	}

	var currentProvider domains.LLMProviderAdapter

	if currentProviderConfig.Compatible == "" {
		return nil, nil
	}

	switch currentProviderConfig.Compatible {
	case "OpenAI":
		currentProvider = providers.NewOpenAICompatibleProviderAdapter(
			currentProviderConfig.BaseURL,
			currentProviderConfig.APIKey,
			currentProviderConfig.AvailableModels,
		)
	}

	agent := GetAgentFromPool()
	agent.WithLLMProvider(currentProvider).
		WithModel(generalConfig.CurrentModel).
		WithState(state).
		WithToolManager(GlobalToolMgr)

	return agent, nil
}

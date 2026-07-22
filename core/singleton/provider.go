package singleton

import (
	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/llm"
	"github.com/smtdfc/nagare/core/llm/providers"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/shared/messages"
)

var AGENT_POOL_SIZE = 10
var GlobalAgentPool *agent.AgentPool
var GlobalConfigMgr *config.ConfigManager
var GlobalConfig *config.Config

var GlobalToolManager *tool.ToolManager

func ensureItemInitialized[T any](item *T) {
	if item == nil {
		panic("Item not initialized")
	}
}

func Init() error {
	var err error
	GlobalAgentPool = agent.NewAgentPool(AGENT_POOL_SIZE).Seed(AGENT_POOL_SIZE)
	GlobalConfigMgr = config.NewConfigManager()
	GlobalConfig, err = GlobalConfigMgr.Load()
	if err != nil {
		return err
	}

	GlobalToolManager = tool.NewToolManager()

	return nil
}

func GetLLMProvider() (llm.LLMProviderAdapter, error) {
	ensureItemInitialized(GlobalConfigMgr)
	ensureItemInitialized(GlobalConfig)

	if GlobalConfig.CurrentProvider == "" {
		return nil, custom_errors.NewLLMProviderError("Provider is not configured. Please check your settings.")
	}

	currentProviderConfig, ok := GlobalConfig.Providers[GlobalConfig.CurrentProvider]
	if !ok {
		return nil, custom_errors.NewLLMProviderError("Provider not found. Please verify the provider name and try again.")
	}

	switch currentProviderConfig.Compatible {
	case config.OPEN_AI:
		return providers.NewOpenAICompatibleProviderAdapter(
			currentProviderConfig.BaseURL,
			currentProviderConfig.APIKey,
			currentProviderConfig.AvailableModels,
		), nil
	}

	return nil, custom_errors.NewLLMProviderError("The selected provider is not compatible with Nagare")
}

func FetchReadyAgent(state *agent.AgentState) (*agent.Agent, error) {
	ensureItemInitialized(GlobalConfigMgr)
	ensureItemInitialized(GlobalConfig)

	provider, err := GetLLMProvider()
	if err != nil {
		return nil, err
	}

	if GlobalConfig.CurrentModel == "" {
		return nil, custom_errors.NewLLMProviderError("No model has been selected. Please choose a model to proceed.")
	}

	agent := GetAgentFromPool()
	agent.WithLLMProvider(provider).
		WithModel(GlobalConfig.CurrentModel).
		WithState(state).
		WithToolManager(GlobalToolManager)

	return agent, nil
}

func CreateEmptyAgentState() *agent.AgentState {
	return agent.NewAgentState(messages.EMPTY_LIST)
}

func GetAgentFromPool() *agent.Agent {
	ensureItemInitialized(GlobalAgentPool)
	return GlobalAgentPool.Get()
}

func PutAgentIntoPool(agent *agent.Agent) {
	ensureItemInitialized(GlobalAgentPool)
	GlobalAgentPool.Put(agent)
}

func CreateContext() *context.ExecuteContext {
	ensureItemInitialized(GlobalToolManager)
	return context.NewExecuteContext(GlobalToolManager)
}

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
	GlobalAgentPool = agent.NewAgentPool(AGENT_POOL_SIZE)
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
		return nil, custom_errors.NewLLMProviderError("Provider not config")
	}

	currentProviderConfig, ok := GlobalConfig.Providers[GlobalConfig.CurrentProvider]
	if !ok {
		return nil, custom_errors.NewLLMProviderError("Provider not found")
	}

	switch currentProviderConfig.Compatible {
	case config.OPEN_AI:
		return providers.NewOpenAICompatibleProviderAdapter(
			currentProviderConfig.BaseURL,
			currentProviderConfig.APIKey,
		), nil
	}

	return nil, custom_errors.NewLLMProviderError("Provider not compatible with nagare")
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

package global

import (
	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/config"
	llm_manager "github.com/smtdfc/nagare/core/llm/manager"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/persistence/database"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/shared/helpers"
)

type Global struct {
}

func Init() error {
	logger.Logger.Info("Global init")
	_, err := database.InitDatabase()
	if err != nil {
		panic(err)
	}

	GlobalLlmRepository = repositories.NewLLMProviderRepository()
	GlobalKVRepository = repositories.NewKVRepository()
	GlobalSessionRepository = repositories.NewSessionRepository()
	GlobalMessageRepository = repositories.NewMessageRepository()
	GlobalPluginRepository = repositories.NewPluginRepository()
	GlobalConfigMgr = config.NewConfigManager(GlobalLlmRepository, GlobalKVRepository)
	GlobalSessionMgr = session.NewSessionManager(GlobalSessionRepository, GlobalMessageRepository)
	GlobalAgentPool = agent.NewAgentPool(AGENT_POOL_SIZE).Seed(AGENT_POOL_SIZE)
	GlobalToolMgr = tool.NewToolManager()
	GlobalLLMManager = llm_manager.NewLLMManager()
	GlobalConnectCodeMgr = plugin.NewConnectCodeManager()
	GlobalPluginMgr = plugin.NewPluginManager(GlobalPluginRepository, GlobalConnectCodeMgr)

	if !helpers.IsRunWithRemoteServer() {
		err = GlobalPluginMgr.Start()
		if err != nil {
			logger.Logger.Error("Failed to start plugins", "error", err)
		}
	}

	return nil
}

package global

import (
	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/persistence/database"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/session"
	"github.com/smtdfc/nagare/core/tool"
)

func Init() error {
	_, err := database.InitDatabase()
	if err != nil {
		panic(err)
	}

	GlobalLlmRepository = repositories.NewLLMProviderRepository()
	GlobalKVRepository = repositories.NewKVRepository()
	GlobalSessionRepository = repositories.NewSessionRepository()
	GlobalMessageRepository = repositories.NewMessageRepository()
	GlobalConfigMgr = config.NewConfigManager(GlobalLlmRepository, GlobalKVRepository)
	GlobalSessionMgr = session.NewSessionManager(GlobalSessionRepository, GlobalMessageRepository)
	GlobalAgentPool = agent.NewAgentPool(AGENT_POOL_SIZE).Seed(AGENT_POOL_SIZE)
	GlobalToolMgr = tool.NewToolManager()
	return nil
}

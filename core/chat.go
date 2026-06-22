package core

import (
	"fmt"

	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/model"
	tool_declarations "github.com/smtdfc/nagare/core/tool/declarations"
)

func InitAgent(conf *config.Config) (*agent.AgentPool, *agent.SessionManager) {
	sessionMgr := agent.NewSessionManager()
	var chatModel model.ChatModel

	currentProvider, ok := conf.Providers[conf.CurrentProvider]
	if !ok || !currentProvider.Enabled {
		return nil, sessionMgr
	}

	if currentProvider.Compatible == config.OPEN_AI_COMPATIBLE {
		chatModel = model.NewOpenAICompatibleClient(&model.ChatModelConfig{
			BaseURL: currentProvider.BaseURL,
			APIKey:  currentProvider.APIKey,
			Model:   currentProvider.ModelName,
		})
	} else {
		fmt.Println("no compatible")
	}

	return agent.NewAgentPool(10, chatModel, tool_declarations.Tools), sessionMgr
}

package core

import (
	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/utils"
)

var SessionMgr *agent.SessionManager
var PluginMgr *plugin.PluginManager
var Config *config.Config
var AgentPool *agent.AgentPool

func SetupEnvironment() error {
	utils.InitGlobalPath()
	err := logger.InitLogger()
	if err != nil {
		return err
	}

	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	Config = config
	SessionMgr = agent.NewSessionManager()
	PluginMgr = plugin.NewPluginManager(config)
	AgentPool = InitAgent(config)
	return nil
}

func PreStart() error {
	go PluginMgr.StartHost()
	return PluginMgr.LoadPlugin()
}

func Shutdown() {
	PluginMgr.Shutdown()
	config.SaveConfig(Config)
}

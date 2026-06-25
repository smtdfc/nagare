package core

import (
	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/tool/declarations"
	nagare_logger "github.com/smtdfc/nagare/shared/logger"
	nagare_path "github.com/smtdfc/nagare/shared/path"
)

var SessionMgr *agent.SessionManager
var PluginMgr *plugin.PluginManager
var Config *config.Config
var AgentPool *agent.AgentPool

func SetupEnvironment() error {
	err := nagare_logger.InitLogger(nagare_path.LogDir)
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
	declarations.InitTools()
	return nil
}

func PreStart() error {
	go PluginMgr.StartHost(AgentPool, SessionMgr)
	return PluginMgr.LoadPlugin()
}

func Shutdown() {
	PluginMgr.Shutdown()
	config.SaveConfig(Config)
}

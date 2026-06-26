package core

import (
	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/plugin/features"
	plugin_host "github.com/smtdfc/nagare/core/plugin/host"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/tool/declarations"
	"github.com/smtdfc/nagare/plugin-sdk/host"
	nagare_logger "github.com/smtdfc/nagare/shared/logger"
	nagare_path "github.com/smtdfc/nagare/shared/path"
)

var GlobalSessionMgr *agent.SessionManager
var GlobalPluginMgr *manager.PluginManager
var GlobalConfig *config.Config
var GlobalAgentPool *agent.AgentPool
var GlobalPluginHost *host.Host
var GlobalChatChannelMgr *features.ChatChannelManager
var GlobalToolRegistry *tool.ToolRegistry

func SetupEnvironment() error {
	err := nagare_logger.InitLogger(nagare_path.LogDir)
	if err != nil {
		return err
	}

	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	GlobalConfig = config
	GlobalToolRegistry = tool.NewToolRegistry()
	declarations.InitTools(GlobalToolRegistry)
	GlobalSessionMgr = agent.NewSessionManager()
	GlobalPluginMgr = manager.NewPluginManager(GlobalConfig)
	GlobalAgentPool = InitAgent(GlobalConfig, GlobalToolRegistry)
	return nil
}

func PreStart() error {
	if GlobalAgentPool == nil || GlobalSessionMgr == nil || GlobalPluginMgr == nil {
		panic("Error: Please call SetupEnvironment first.")
	}

	GlobalPluginHost = host.NewHost(nagare_logger.GetLogger("Plugin Host"))
	GlobalChatChannelMgr = features.NewChatChannelManager(GlobalPluginHost)
	go plugin_host.StartHost(GlobalPluginHost, GlobalPluginMgr, GlobalChatChannelMgr, GlobalAgentPool, GlobalSessionMgr)
	return GlobalPluginMgr.LoadPlugin()
}

func Shutdown() {
	GlobalPluginMgr.Shutdown()
	config.SaveConfig(GlobalConfig)
}

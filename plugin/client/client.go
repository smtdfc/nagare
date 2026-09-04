package client

import (
	"os"

	"github.com/smtdfc/nagare/plugin/metadata"
)

type PluginClient struct {
	Metadata *metadata.PluginMetadata
	Config   *Config
}

func (p *PluginClient) Start() {
	p.Config.Port = os.Getenv("NAGARE_PLUGIN_HOST_PORT")
	p.Config.ConnectCode = os.Getenv("NAGARE_PLUGIN_CONNECT_CODE")
}

func NewPlugin() *PluginClient {
	return &PluginClient{
		Config: &Config{},
	}
}

package plugin_sdk

import (
	"github.com/hashicorp/go-plugin"
	"github.com/smtdfc/nagare/plugin-sdk/proto"
)

const PLUGIN_MAGIC_COOKIE_KEY = "NAGARE_PLUGIN"
const PLUGIN_MAGIC_COOKIE_VALUE = "hello"
const PLUGIN_PROTOCOL_VERSION = 1

type PluginBase struct {
	proto.UnimplementedPluginServer
}

func StartServer(impl proto.PluginServer) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  PLUGIN_PROTOCOL_VERSION,
			MagicCookieKey:   PLUGIN_MAGIC_COOKIE_KEY,
			MagicCookieValue: PLUGIN_MAGIC_COOKIE_VALUE,
		},
		Plugins: map[string]plugin.Plugin{
			"plugin": &NagarePlugin{Impl: impl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

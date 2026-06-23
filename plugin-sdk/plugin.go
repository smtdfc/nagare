package plugin_sdk

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/smtdfc/nagare/plugin-sdk/proto"
	"google.golang.org/grpc"
)

type NagarePlugin struct {
	plugin.Plugin
	Impl proto.PluginServer
}

func (p *NagarePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPluginServer(s, p.Impl)
	return nil
}

func (p *NagarePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return proto.NewPluginClient(c), nil
}

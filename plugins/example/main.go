package main

import (
	"context"

	plugin_sdk "github.com/smtdfc/nagare/plugin-sdk"
	"github.com/smtdfc/nagare/plugin-sdk/proto"
)

type MyPlugin struct {
	plugin_sdk.PluginBase
}

func (p *MyPlugin) Metadata(ctx context.Context, req *proto.Empty) (*proto.PluginMetadata, error) {
	return &proto.PluginMetadata{
		Id:   "org.smtdfc.nagare.plugins.example",
		Name: "Hello World Plugin",
	}, nil
}

func main() {
	plugin_sdk.StartServer(&MyPlugin{})
}

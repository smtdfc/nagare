package main

import (
	"context"
	_ "embed"

	"github.com/smtdfc/nagare/plugin/client"
)

//go:embed metadata.json
var metadata string

func main() {
	ctx := context.Background()
	pluginClient := client.NewPlugin()
	_, err := pluginClient.LoadMetadata(metadata)
	if err != nil {
		pluginClient.Logger.Error("Load metadata error", "error", err)
		return
	}

	err = pluginClient.Start(
		ctx,
		func() {
			err := pluginClient.Handshake(ctx)
			if err != nil {
				pluginClient.Logger.Error("Handshake error", "error", err)
				return
			}
		},
	)
	if err != nil {
		pluginClient.Logger.Error("Handshake error", "error", err)
	}
}

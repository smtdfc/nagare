package main

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/smtdfc/nagare/plugin/client"
)

//go:embed metadata.json
var metadata string

func main() {
	ctx := context.Background()
	pluginClient := client.NewPlugin()
	_, err := pluginClient.LoadMetadata(metadata)
	if err != nil {
		fmt.Println("Metadata error :", err)
		return
	}

	err = pluginClient.Start(
		ctx,
		func() {
			err := pluginClient.Handshake(ctx)
			if err != nil {
				fmt.Println("Handshake error :", err)
				return
			}
		},
	)
	if err != nil {
		fmt.Println("Plugin start error :", err)
	}
}

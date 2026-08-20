package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/smtdfc/nagare/plugin/client"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client := client.NewPluginClient("example")
	go client.Start(ctx, OnReady)

	<-ctx.Done()
}

func OnReady(ctx context.Context) error {
	return nil
}

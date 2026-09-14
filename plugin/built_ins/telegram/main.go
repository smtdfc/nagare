package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/mymmrac/telego"
	"github.com/smtdfc/nagare/plugin/client"
)

//go:embed metadata.json
var metadata string

func main() {
	token := os.Getenv("NAGARE_TELEGRAM_TOKEN")
	ctx := context.Background()
	bot, err := telego.NewBot(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pluginClient := client.NewPlugin()
	telegramPlugin := NewTelegramPlugin(bot, pluginClient)

	err = telegramPlugin.Start(ctx)
	if err != nil {
		fmt.Println("Plugin start error :", err)
	}
}

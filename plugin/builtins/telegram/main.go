package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/smtdfc/nagare/plugin/client"
)

type PluginConfig struct {
	Token string `json:"token"`
}

var plugin = client.NewPluginClient("nagare.telegram")
var telegramBot *bot.Bot

func main() {
	var err error
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config, err := client.ReadPluginConfig[PluginConfig](plugin)
	if err != nil {
		plugin.Logger.Error("Failed to read config", "error", err)
		return
	}

	if config == nil || config.Token == "" {
		plugin.Logger.Error("Telegram bot token is missing in config file")
		return
	}

	telegramBot, err = bot.New(config.Token, bot.WithDefaultHandler(handler))
	if err != nil {
		plugin.Logger.Error("Failed to init Telegram bot", "error", err)
		return
	}

	go func() {
		err := plugin.Start(ctx, OnReady)
		if err != nil {
			plugin.Logger.Error("Failed to start plugin", "error", err)
			return
		}
	}()

	<-ctx.Done()
}

func OnReady(ctx context.Context) error {
	plugin.Logger.Debug("Plugin ready")
	err := plugin.Register()
	if err != nil {
		plugin.Logger.Error("Failed to register", "error", err)
		return err
	}

	plugin.Logger.Debug("Starting Telegram bot")
	go telegramBot.Start(ctx)
	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := update.Message.Text
	sender := update.Message.From.FirstName

	plugin.Logger.Debug("Recived message", "sender", sender, "text", text)
}

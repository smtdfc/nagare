package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
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

	go func() {
		err := plugin.Start(ctx, func(ctx context.Context) error {
			plugin.Logger.Debug("Plugin ready")
			err := plugin.Register()
			if err != nil {
				plugin.Logger.Error("Failed to register", "error", err)
				return err
			}

			plugin.Logger.Debug("Starting Telegram bot")
			InitBot(config, ctx)
			return nil
		})
		if err != nil {
			plugin.Logger.Error("Failed to start plugin", "error", err)
			return
		}
	}()

	<-ctx.Done()
}

func InitBot(config *PluginConfig, ctx context.Context) {
	var telegramBot *bot.Bot
	var err error
	var botHandler = NewBotHandler()
	customClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.handler),
		bot.WithHTTPClient(30*time.Minute, customClient),
	}
	backoff := 2 * time.Second
	maxBackoff := 30 * time.Second

	for {
		plugin.Logger.Info("Attempting to initialize Telegram bot...")

		telegramBot, err = bot.New(config.Token, opts...)
		if err == nil {
			plugin.Logger.Info("Telegram bot initialized successfully!")
			break
		}

		plugin.Logger.Error("Failed to init Telegram bot, retrying...", "error", err, "retry_in", backoff)
		time.Sleep(backoff)

		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}

	go telegramBot.Start(ctx)
}

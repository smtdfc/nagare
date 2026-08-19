package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

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

var (
	activeSenders = make(map[int64]bool)
	queueMu       sync.Mutex
)

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID
	userName := update.Message.From.FirstName
	text := update.Message.Text
	chatID := update.Message.Chat.ID

	queueMu.Lock()
	if isActive, ok := activeSenders[userID]; ok && isActive {
		queueMu.Unlock()
		plugin.Logger.Debug("Rejected message because previous task is running", "user_id", userID, "sender", userName, "text", text)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Please wait for the previous request to finish processing!",
		})
		if err != nil {
			plugin.Logger.Error("Failed to send pending notice", "error", err)
		}
		return
	}

	activeSenders[userID] = true
	queueMu.Unlock()

	plugin.Logger.Debug("Received message", "user_id", userID, "sender", userName, "text", text)

	// TODO: Xử lý logic nặng ở đây...
}

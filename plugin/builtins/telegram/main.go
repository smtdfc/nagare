package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/smtdfc/nagare/plugin/client"
	"github.com/smtdfc/nagare/shared/messages"
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
	customClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
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
	defer func() {
		queueMu.Lock()
		delete(activeSenders, userID)
		queueMu.Unlock()
	}()

	plugin.Logger.Debug("Received message", "user_id", userID, "sender", userName, "text", text)
	output, err := plugin.InvokeAgent(text, fmt.Sprintf("%d", userID))
	if err != nil {
		plugin.Logger.Error("Failed to invoke agent", "error", err)
		return
	}

	var textBuffer strings.Builder

	for chunk := range output {
		switch message := chunk.(type) {
		case *messages.Text:
			textBuffer.WriteString(message.Content)

		case *messages.ToolCall:
			if textBuffer.Len() > 0 {
				if content := strings.TrimSpace(textBuffer.String()); content != "" {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatID,
						Text:   content,
					})
				}
				textBuffer.Reset()
			}

			if strings.TrimSpace(message.Name) != "" {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   fmt.Sprintf("Call: %s", message.Name),
				})
			}
		}
	}

	if finalText := strings.TrimSpace(textBuffer.String()); finalText != "" {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   finalText,
		})
		if err != nil {
			plugin.Logger.Error("Failed to send text message", "error", err)
		}
	}
}

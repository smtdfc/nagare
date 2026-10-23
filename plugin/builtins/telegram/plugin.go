package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/smtdfc/nagare/plugin/client"
)

type TelegramPlugin struct {
	mu           sync.Mutex
	bot          *telego.Bot
	pluginClient *client.PluginClient
	channels     map[string]string
	idleTimers   map[string]*time.Timer

	queues         map[string][]telego.Update
	isProcessing   map[string]bool
	messageBuffers sync.Map
	approvals      map[string]pendingToolApproval
}

type pendingToolApproval struct {
	chatID    string
	sessionID string
	messageID int
}

func NewTelegramPlugin(bot *telego.Bot, pluginClient *client.PluginClient) *TelegramPlugin {
	return &TelegramPlugin{
		bot:          bot,
		pluginClient: pluginClient,
		channels:     make(map[string]string),
		idleTimers:   make(map[string]*time.Timer),
		queues:       make(map[string][]telego.Update),
		isProcessing: make(map[string]bool),
		approvals:    make(map[string]pendingToolApproval),
	}
}

func (tp *TelegramPlugin) Start(ctx context.Context) error {
	_, err := tp.pluginClient.LoadMetadata(metadata)
	if err != nil {
		return err
	}

	tp.pluginClient.OnReceivedChatMessage = tp.OnReceivedChatMessage
	return tp.pluginClient.Start(
		ctx,
		func() {
			err := tp.pluginClient.Handshake(ctx)
			if err != nil {
				return
			}

			updates, _ := tp.bot.UpdatesViaLongPolling(ctx, nil)
			for update := range updates {
				if update.CallbackQuery != nil {
					tp.handleCallbackQuery(update)
					continue
				}
				if update.Message != nil {
					tp.enqueueMessage(update)
				}
			}
		},
	)
}

// approvalKeyboard builds the inline keyboard with Approve / Reject buttons
// for a pending tool approval request.
func (tp *TelegramPlugin) approvalKeyboard(requestID string) *telego.InlineKeyboardMarkup {
	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{
				{Text: "✅ Approve", CallbackData: fmt.Sprintf("approve:%s", requestID)},
				{Text: "❌ Reject", CallbackData: fmt.Sprintf("reject:%s", requestID)},
			},
		},
	}
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

// sendTextMessage sends a text messages to a Telegram chat and logs any error.
func (tp *TelegramPlugin) sendTextMessage(ctx context.Context, chatID int64, text string) error {
	if text == "" {
		return nil
	}
	params := tu.Message(tu.ID(chatID), text)
	_, err := tp.bot.SendMessage(ctx, params)
	if err != nil {
		tp.pluginClient.Logger.Error("Error sending messages to telegram", "chatID", chatID, "error", err)
	}
	return err
}

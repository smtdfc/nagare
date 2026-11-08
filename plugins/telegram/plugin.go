package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	sdk "github.com/smtdfc/nagare/pkgs/plugin/sdk"
)

// discardLogger is used when the plugin client is not wired up, so unit tests
// can construct a TelegramPlugin without a logger.
var discardLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

// logger returns the plugin logger, falling back to a no-op logger when the
// plugin client is missing.
func (tp *TelegramPlugin) logger() *slog.Logger {
	if tp.pluginClient != nil && tp.pluginClient.Logger != nil {
		return tp.pluginClient.Logger
	}
	return discardLogger
}

type TelegramPlugin struct {
	mu           sync.Mutex
	bot          *telego.Bot
	pluginClient *sdk.PluginClient
	channels     map[string]struct{}
	idleTimers   map[string]*time.Timer

	queues         map[string][]telego.Update
	isProcessing   map[string]bool
	messageBuffers sync.Map
	approvals      map[string]pendingToolApproval
}

type pendingToolApproval struct {
	channelID string
	messageID int
}

func NewTelegramPlugin(bot *telego.Bot, pluginClient *sdk.PluginClient) *TelegramPlugin {
	return &TelegramPlugin{
		bot:          bot,
		pluginClient: pluginClient,
		channels:     make(map[string]struct{}),
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

// sendTextMessage sends a text message to a Telegram chat, splitting it into
// multiple messages when it exceeds Telegram's message length limit.
func (tp *TelegramPlugin) sendTextMessage(ctx context.Context, chatID int64, text string) error {
	if strings.TrimSpace(text) == "" {
		return nil
	}

	for _, chunk := range splitMessage(text, telegramMaxMessageLength) {
		params := tu.Message(tu.ID(chatID), chunk)
		if _, err := tp.bot.SendMessage(ctx, params); err != nil {
			tp.logger().Error("Error sending messages to telegram", "chatID", chatID, "error", err)
			return err
		}
	}

	return nil
}

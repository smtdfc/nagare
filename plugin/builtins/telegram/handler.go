package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/smtdfc/nagare/shared/messages"
)

type BotHandler struct {
	activeSenders map[int64]bool
	queueMu       sync.Mutex
}

func (h *BotHandler) InvokeAgent(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	text := update.Message.Text
	chatID := update.Message.Chat.ID

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

func (h *BotHandler) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID
	userName := update.Message.From.FirstName
	text := update.Message.Text
	chatID := update.Message.Chat.ID

	h.queueMu.Lock()
	if isActive, ok := h.activeSenders[userID]; ok && isActive {
		h.queueMu.Unlock()
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

	h.activeSenders[userID] = true
	h.queueMu.Unlock()
	defer func() {
		h.queueMu.Lock()
		delete(h.activeSenders, userID)
		h.queueMu.Unlock()
	}()

	if len(text) > 0 && text[0] == '/' {
		h.handleCommand(ctx, b, update, text)
		return
	}

	h.InvokeAgent(ctx, b, update)
}

func (h *BotHandler) handleCommand(ctx context.Context, b *bot.Bot, update *models.Update, text string) {
	parsed := ParseCommand(text)
	if parsed == nil {
		return
	}

	plugin.Logger.Debug("Parsed command", "cmd", parsed.Command, "args", parsed.Args)
	switch parsed.Command {
	case "start":
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Nagare is Ready",
		})

	case "reset":
		h.ResetChat(ctx, b, update)
	default:
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Unknown command. Please try again!",
		})
	}
}

func (h *BotHandler) ResetChat(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	sessionID := fmt.Sprintf("%d", userID)
	chatID := update.Message.Chat.ID
	sentMsg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "🧹 Dusting off the cobwebs of our memory, please hold on...",
	})

	if err != nil {
		plugin.Logger.Error("Failed to send initial message", "error", err)
		return
	}

	err = plugin.ResetChatSession(sessionID)
	var responseText string
	if err != nil {
		plugin.Logger.Error("Failed to reset chat session", "error", err, "user_id", userID)
		responseText = "⚠️ Oops! My brain got stuck and I couldn't wipe the memory: " + err.Error()
	} else {
		responseText = "✨ All done! History wiped clean. Let's start with a fresh slate, my liege."
	}

	_, _ = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatID,
		MessageID: sentMsg.ID,
		Text:      responseText,
	})
}

func NewBotHandler() *BotHandler {
	return &BotHandler{
		activeSenders: map[int64]bool{},
		queueMu:       sync.Mutex{},
	}
}

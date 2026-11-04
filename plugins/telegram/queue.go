package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (tp *TelegramPlugin) enqueueMessage(update telego.Update) {
	channelID := strconv.FormatInt(update.Message.Chat.ID, 10)

	tp.logger().Info("enqueueMessage", "channelID", channelID)
	tp.mu.Lock()
	tp.resetIdleTimer(channelID)
	tp.queues[channelID] = append(tp.queues[channelID], update)

	if !tp.isProcessing[channelID] {
		tp.isProcessing[channelID] = true
		tp.logger().Info("processing", "channelID", channelID)
		go tp.processQueue(context.Background(), channelID)
	}
	tp.mu.Unlock()
}

func (tp *TelegramPlugin) processQueue(ctx context.Context, channelID string) {
	for {
		tp.mu.Lock()
		q := tp.queues[channelID]
		if len(q) == 0 {
			tp.isProcessing[channelID] = false
			tp.mu.Unlock()
			return
		}

		update := q[0]
		tp.queues[channelID] = q[1:]
		tp.mu.Unlock()

		tp.handleSingleMessage(ctx, &update)
	}
}

func (tp *TelegramPlugin) handleSingleMessage(ctx context.Context, update *telego.Update) {
	chatIntID := update.Message.Chat.ID
	channelID := strconv.FormatInt(chatIntID, 10)

	if !tp.hasChannel(channelID) {
		if err := tp.pluginClient.PrepareChatChannel(ctx, channelID); err != nil {
			_ = tp.sendTextMessage(ctx, chatIntID, "Oops! Error while preparing chat channel")
			return
		}
		tp.setChannel(channelID)
	}

	if tp.handleApprovalCommand(ctx, channelID, chatIntID, update.Message.Text) {
		return
	}

	switch update.Message.Text {
	case "/start":
		_ = tp.sendTextMessage(ctx, chatIntID, "Hello! I'm your Nagare bot. How can I assist you today?")
		tp.finishProcessing(channelID)
		return
	case "/help":
		_ = tp.sendTextMessage(ctx, chatIntID, "You can send me any messages and I'll process it for you.")
		tp.finishProcessing(channelID)
		return
	case "/ping":
		_ = tp.sendTextMessage(ctx, chatIntID, "Pong!")
		tp.finishProcessing(channelID)
		return
	case "/chat_id":
		_ = tp.sendTextMessage(ctx, chatIntID, "Your channel ID is: "+channelID)
		tp.finishProcessing(channelID)
		return
	case "/reset":
		err := tp.pluginClient.ResetChatChannel(ctx, channelID)
		if err != nil {
			tp.logger().Error("failed to reset chat channel", "channelID", channelID, "error", err)
			_ = tp.sendTextMessage(ctx, chatIntID, "Oops! Error while resetting chat channel")
			tp.finishProcessing(channelID)
			return
		}
		tp.clearChannel(channelID)
		_ = tp.sendTextMessage(ctx, chatIntID, "Chat channel has been reset.")
		tp.finishProcessing(channelID)
		return

	case "/model_providers":
		providers, err := tp.pluginClient.GetListModelProvider(ctx)
		if err != nil {
			tp.logger().Error("failed to get list model provider", "channelID", channelID, "error", err)
			_ = tp.sendTextMessage(ctx, chatIntID, "Oops! Error while getting list model providers")
			tp.finishProcessing(channelID)
			return
		}

		var msg strings.Builder
		msg.WriteString("List of model providers:\n")
		idx := 0
		for id, p := range providers {
			msg.WriteString(fmt.Sprintf("%d: %s (ID: `%s`)\n", idx, p, id))
			idx = idx + 1
		}
		_ = tp.sendTextMessage(ctx, chatIntID, msg.String())
		tp.finishProcessing(channelID)
		return
	}

	err := tp.pluginClient.SendChatMessage(ctx, channelID, update.Message.Text)
	if err != nil {
		_ = tp.sendTextMessage(ctx, chatIntID, "Oops! Error while sending messages")
		tp.finishProcessing(channelID)
		return
	}
}

func (tp *TelegramPlugin) handleApprovalCommand(ctx context.Context, channelID string, chatIntID int64, text string) bool {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return false
	}
	command := strings.ToLower(strings.Split(fields[0], "@")[0])
	if command != "/approve" && command != "/reject" {
		return false
	}
	if len(fields) != 2 {
		_ = tp.sendTextMessage(ctx, chatIntID, "Usage: /approve <request-id> or /reject <request-id>")
		return true
	}

	approval, ok := tp.getPendingApproval(fields[1], channelID)
	if !ok {
		_ = tp.sendTextMessage(ctx, chatIntID, "That approval request is no longer pending.")
		return true
	}
	accept := command == "/approve"
	if err := tp.pluginClient.SubmitToolApproval(ctx, approval.channelID, fields[1], accept); err != nil {
		tp.logger().Error("failed to submit tool approval", "requestID", fields[1], "error", err)
		_ = tp.sendTextMessage(ctx, chatIntID, "Could not submit your decision. Please try again.")
		return true
	}
	tp.clearPendingApproval(fields[1])
	if accept {
		_ = tp.sendTextMessage(ctx, chatIntID, "✅ Tool execution approved.")
	} else {
		_ = tp.sendTextMessage(ctx, chatIntID, "❌ Tool execution rejected.")
	}
	return true
}

func (tp *TelegramPlugin) handleCallbackQuery(update telego.Update) {
	cq := update.CallbackQuery
	if cq == nil || cq.Message == nil {
		return
	}

	chatIntID := cq.Message.GetChat().ID
	channelID := strconv.FormatInt(chatIntID, 10)

	parts := strings.SplitN(cq.Data, ":", 2)
	if len(parts) != 2 {
		_ = tp.bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Invalid action.",
		})
		return
	}

	action := parts[0]
	requestID := parts[1]

	var accept bool
	switch action {
	case "approve":
		accept = true
	case "reject":
		accept = false
	default:
		_ = tp.bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Invalid action.",
		})
		return
	}

	approval, ok := tp.getPendingApproval(requestID, channelID)
	if !ok {
		_ = tp.bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "That approval request is no longer pending.",
		})
		return
	}

	if err := tp.pluginClient.SubmitToolApproval(context.Background(), approval.channelID, requestID, accept); err != nil {
		tp.logger().Error("failed to submit tool approval via button", "requestID", requestID, "error", err)
		_ = tp.bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Could not submit your decision. Please try again.",
		})
		return
	}

	// Acknowledge the callback to remove the loading state on the button.
	_ = tp.bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cq.ID,
	})

	// Remove the inline keyboard from the original request message.
	_, _ = tp.bot.EditMessageReplyMarkup(context.Background(), &telego.EditMessageReplyMarkupParams{
		ChatID:      tu.ID(chatIntID),
		MessageID:   approval.messageID,
		ReplyMarkup: &telego.InlineKeyboardMarkup{},
	})

	if accept {
		_ = tp.sendTextMessage(context.Background(), chatIntID, "✅ Tool execution approved.")
	} else {
		_ = tp.sendTextMessage(context.Background(), chatIntID, "❌ Tool execution rejected.")
	}

	tp.clearPendingApproval(requestID)
}

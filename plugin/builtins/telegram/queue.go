package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
)

func (tp *TelegramPlugin) enqueueMessage(update telego.Update) {
	chatID := strconv.FormatInt(update.Message.Chat.ID, 10)

	tp.pluginClient.Logger.Info("enqueueMessage", "chatID", chatID)
	tp.mu.Lock()
	tp.resetIdleTimer(chatID)
	tp.queues[chatID] = append(tp.queues[chatID], update)

	if !tp.isProcessing[chatID] {
		tp.isProcessing[chatID] = true
		tp.pluginClient.Logger.Info("processing", "chatID", chatID)
		go tp.processQueue(context.Background(), chatID)
	}
	tp.mu.Unlock()
}

func (tp *TelegramPlugin) processQueue(ctx context.Context, chatID string) {
	for {
		tp.mu.Lock()
		q := tp.queues[chatID]
		if len(q) == 0 {
			tp.isProcessing[chatID] = false
			tp.mu.Unlock()
			return
		}

		update := q[0]
		tp.queues[chatID] = q[1:]
		tp.mu.Unlock()

		tp.handleSingleMessage(ctx, &update)
	}
}

func (tp *TelegramPlugin) clearSessionID(chatID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	delete(tp.channels, chatID)
}

func (tp *TelegramPlugin) handleSingleMessage(ctx context.Context, update *telego.Update) {
	chatIntID := update.Message.Chat.ID
	chatID := strconv.FormatInt(chatIntID, 10)
	sessionID, exist := tp.getSessionID(chatID)
	if !exist {
		var err error
		sessionID, err = tp.pluginClient.PrepareChatSession(ctx, chatID)
		if err != nil {
			_ = tp.sendTextMessage(ctx, chatIntID, "Oops! Error while preparing chat session")
			return
		}
		tp.setSessionID(chatID, sessionID)
	}

	if tp.handleApprovalCommand(ctx, chatID, chatIntID, update.Message.Text) {
		return
	}

	switch update.Message.Text {
	case "/start":
		_ = tp.sendTextMessage(ctx, chatIntID, "Hello! I'm your Nagare bot. How can I assist you today?")
		tp.finishProcessing(chatID, sessionID)
		return
	case "/help":
		_ = tp.sendTextMessage(ctx, chatIntID, "You can send me any messages and I'll process it for you.")
		tp.finishProcessing(chatID, sessionID)
		return
	case "/ping":
		_ = tp.sendTextMessage(ctx, chatIntID, "Pong!")
		tp.finishProcessing(chatID, sessionID)
		return
	case "/chat_id":
		_ = tp.sendTextMessage(ctx, chatIntID, "Your chat ID is: "+chatID)
		tp.finishProcessing(chatID, sessionID)
		return
	case "/reset":
		err := tp.pluginClient.ResetChatChannel(ctx, chatID)
		if err != nil {
			tp.pluginClient.Logger.Error("failed to reset chat session", "chatID", chatID, "error", err)
			_ = tp.sendTextMessage(ctx, chatIntID, "Oops! Error while resetting chat session")
			tp.finishProcessing(chatID, sessionID)
			return
		}
		tp.clearSessionID(chatID)
		_ = tp.sendTextMessage(ctx, chatIntID, "Chat session has been reset.")
		tp.finishProcessing(chatID, sessionID)
		return

	case "/model_providers":
		providers, err := tp.pluginClient.GetListModelProvider(ctx)
		if err != nil {
			tp.pluginClient.Logger.Error("failed to get list model provider", "chatID", chatID, "error", err)
			_ = tp.sendTextMessage(ctx, chatIntID, "Oops! Error while getting list model providers")
			tp.finishProcessing(chatID, sessionID)
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
		tp.finishProcessing(chatID, sessionID)
		return
	}

	err := tp.pluginClient.SendChatMessage(ctx, sessionID, update.Message.Text)
	if err != nil {
		_ = tp.sendTextMessage(ctx, chatIntID, "Oops! Error while sending messages")
		tp.finishProcessing(chatID, sessionID)
		return
	}
}

func (tp *TelegramPlugin) handleApprovalCommand(ctx context.Context, chatID string, chatIntID int64, text string) bool {
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

	approval, ok := tp.getPendingApproval(fields[1], chatID)
	if !ok {
		_ = tp.sendTextMessage(ctx, chatIntID, "That approval request is no longer pending.")
		return true
	}
	accept := command == "/approve"
	if err := tp.pluginClient.SubmitToolApproval(ctx, approval.sessionID, fields[1], accept); err != nil {
		tp.pluginClient.Logger.Error("failed to submit tool approval", "requestID", fields[1], "error", err)
		_ = tp.sendTextMessage(ctx, chatIntID, "Could not submit your decision. Please try again.")
		return true
	}
	tp.clearPendingApproval(fields[1])
	if accept {
		_ = tp.sendTextMessage(ctx, chatIntID, "Tool execution approved.")
	} else {
		_ = tp.sendTextMessage(ctx, chatIntID, "Tool execution rejected.")
	}
	return true
}

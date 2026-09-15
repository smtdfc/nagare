package main

import (
	"context"
	"strconv"

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

func (tp *TelegramPlugin) handleSingleMessage(ctx context.Context, update *telego.Update) {
	chatIntID := update.Message.Chat.ID
	chatID := strconv.FormatInt(chatIntID, 10)

	switch update.Message.Text {
	case "/start":
		_ = tp.sendTextMessage(ctx, chatIntID, "Hello! I'm your Nagare bot. How can I assist you today?")
		return
	case "/help":
		_ = tp.sendTextMessage(ctx, chatIntID, "You can send me any message and I'll process it for you.")
		return
	case "/ping":
		_ = tp.sendTextMessage(ctx, chatIntID, "Pong!")
		return
	case "/chat_id":
		_ = tp.sendTextMessage(ctx, chatIntID, "Your chat ID is: "+chatID)
		return
	}

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

	err := tp.pluginClient.SendChatMessage(ctx, sessionID, update.Message.Text)
	if err != nil {
		_ = tp.sendTextMessage(ctx, chatIntID, "Oops! Error while sending message")
		return
	}
}

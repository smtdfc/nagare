package main

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/smtdfc/nagare/plugin-sdk/plugin"
	"github.com/smtdfc/nagare/plugin-sdk/shared"
)

func HandlePluginMessages(ctx context.Context, b *bot.Bot, plg *plugin.Plugin, state *BotState, msg shared.Message) {
	switch msg.Kind {
	case shared.SHUTDOWN_PLUGIN_REQUEST:
	case shared.REGISTER_CHAT_CHANNEL_SUCCESS:
		var payload shared.RegisterChatChannelSuccessPayload
		json.Unmarshal(msg.Payload, &payload)
		parts := strings.Split(payload.ID, ":")
		var chatID int64
		if len(parts) >= 2 {
			chatID, _ = strconv.ParseInt(parts[1], 10, 64)
		}

		state.Lock()
		state.isInitialized[chatID] = true
		for _, text := range state.msgQueue[chatID] {
			plg.Send(shared.HANDLE_CHAT_MESSAGE, shared.HandleChatMessagePayload{
				Channel: payload.ID,
				Message: text,
			})
		}
		delete(state.msgQueue, chatID)
		state.Unlock()

	case shared.HANDLE_CHAT_MESSAGE:
		var payload shared.HandleChatMessagePayload
		json.Unmarshal(msg.Payload, &payload)
		parts := strings.Split(payload.Channel, ":")
		var extractedChatID int64
		if len(parts) >= 2 {
			extractedChatID, _ = strconv.ParseInt(parts[1], 10, 64)
		}
		if extractedChatID != 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: extractedChatID,
				Text:   payload.Message,
			})
		}
	}
}

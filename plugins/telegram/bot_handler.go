package main

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/smtdfc/nagare/plugin-sdk/plugin"
	"github.com/smtdfc/nagare/plugin-sdk/shared"
)

func NewTelegramHandler(state *BotState, plg *plugin.Plugin) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil || update.Message.Text == "" {
			return
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		state.Lock()
		if !state.isInitialized[chatID] {
			state.msgQueue[chatID] = append(state.msgQueue[chatID], text)
			if _, exists := state.userChannels[chatID]; !exists {
				channelID := fmt.Sprintf("telegram:%d", chatID)
				state.userChannels[chatID] = channelID
				plg.Send(shared.REGISTER_CHAT_CHANNEL, shared.RegisterChatChannelPayload{ID: channelID})
			}
			state.Unlock()
			return
		}

		channelID := state.userChannels[chatID]
		state.Unlock()

		plg.Send(shared.HANDLE_CHAT_MESSAGE, shared.HandleChatMessagePayload{
			Channel: channelID,
			Message: text,
		})
	}
}

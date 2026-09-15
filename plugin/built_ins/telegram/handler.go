package main

import (
	"context"

	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/message"
)

func (tp *TelegramPlugin) OnReceivedChatMessage(sessionID, chunk string) {
	chatID, chatIntID, ok := tp.getChatIDBySessionID(sessionID)
	if !ok {
		return
	}

	msg, err := helpers.UnmarshalJson[message.AnyMessage](chunk)
	if err != nil {
		return
	}

	switch msg.Type {
	case message.TextMessageType:
		textMsg, err := helpers.UnmarshalJson[message.TextMessage](chunk)
		if err != nil {
			return
		}
		buf := tp.getOrCreateBuffer(sessionID)
		buf.Append(textMsg.Content)

	case message.AgentErrorMessageType:
		errMsg, err := helpers.UnmarshalJson[message.AgentErrorMessage](chunk)
		if err != nil {
			return
		}
		tp.pluginClient.Logger.Error("Agent error message received", "chatID", chatID, "sessionID", sessionID, "error", errMsg.Error, "code", errMsg.Code)
		_ = tp.sendTextMessage(context.Background(), chatIntID, errMsg.Error)
		tp.finishProcessing(chatID, sessionID)

	case message.AgentCompletedMessageType:
		completedMsg, err := helpers.UnmarshalJson[message.AgentCompletedMessage](chunk)
		if err != nil {
			return
		}
		if completedMsg.Success {
			buf := tp.getOrCreateBuffer(sessionID)
			_ = tp.sendTextMessage(context.Background(), chatIntID, buf.String())
		}
		tp.finishProcessing(chatID, sessionID)
	}
}

package main

import (
	"context"
	"strconv"

	"github.com/smtdfc/nagare/pkgs/helpers"
	"github.com/smtdfc/nagare/pkgs/messages"
)

func (tp *TelegramPlugin) OnReceivedChatMessage(sessionID, channelID, chunk string) {
	var chatID string
	var chatIntID int64

	if channelID != "" {
		num, err := strconv.ParseInt(channelID, 10, 64)
		if err != nil {
			return
		}
		chatID = channelID
		chatIntID = num
	} else {
		var ok bool
		chatID, chatIntID, ok = tp.getChatIDBySessionID(sessionID)
		if !ok {
			return
		}
	}

	msg, err := helpers.UnmarshalJson[messages.AnyMessage](chunk)
	if err != nil {
		return
	}

	switch msg.Type {
	case messages.ApprovalRequestMessageType:
		approval, err := helpers.UnmarshalJson[messages.ApprovalRequestMessage](chunk)
		if err != nil {
			return
		}
		tp.setPendingApproval(approval.RequestID, chatID, sessionID)

		_ = tp.sendTextMessage(context.Background(), chatIntID,
			"Approval required: .\n"+
				"Reply /approve "+approval.RequestID+" to allow it, or /reject "+approval.RequestID+" to deny it.")

	case messages.TextMessageType:
		textMsg, err := helpers.UnmarshalJson[messages.TextMessage](chunk)
		if err != nil {
			return
		}
		buf := tp.getOrCreateBuffer(sessionID)
		buf.Append(textMsg.Content)

	case messages.AgentErrorMessageType:
		errMsg, err := helpers.UnmarshalJson[messages.AgentErrorMessage](chunk)
		if err != nil {
			return
		}
		tp.pluginClient.Logger.Error("Agent error message received", "chatID", chatID, "sessionID", sessionID, "error", errMsg.Error, "code", errMsg.Code)
		_ = tp.sendTextMessage(context.Background(), chatIntID, errMsg.Error)
		tp.finishProcessing(chatID, sessionID)

	case messages.AgentCompletedMessageType:
		completedMsg, err := helpers.UnmarshalJson[messages.AgentCompletedMessage](chunk)
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

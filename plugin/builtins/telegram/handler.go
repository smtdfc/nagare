package main

import (
	"context"
	"fmt"
	"strconv"

	tu "github.com/mymmrac/telego/telegoutil"
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

		title := escapeMarkdownV2(approval.Title)
		text := fmt.Sprintf("⚠️ *Approval required*\n\nTool: *%s*\n\nDo you allow this tool execution?", title)

		params := tu.Message(tu.ID(chatIntID), text)
		params.ParseMode = "MarkdownV2"
		params.ReplyMarkup = tp.approvalKeyboard(approval.RequestID)

		sent, err := tp.bot.SendMessage(context.Background(), params)
		if err != nil {
			tp.pluginClient.Logger.Error("failed to send approval request", "requestID", approval.RequestID, "chatID", chatID, "error", err)
			return
		}

		tp.setPendingApproval(approval.RequestID, chatID, sessionID, sent.MessageID)

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

package host

import (
	"context"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/event_bus"
	core_plugin "github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/pkgs/messages"
	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func (h *PluginHost) handlePrepareChatChannel(conn *PluginConnection, payload *contracts.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(contracts.PrepareChatChannelFailedEvent, &contracts.PrepareChatChannelFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[contracts.PrepareChatChannelEventPayload](payload)
	if err != nil {
		_ = conn.Send(contracts.PrepareChatChannelFailedEvent, &contracts.PrepareChatChannelFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !RequireScopes(conn.scopes, core_plugin.ChatFeature) {
		_ = conn.Send(contracts.PrepareChatChannelFailedEvent, &contracts.PrepareChatChannelFailedEventPayload{
			ChannelID: data.ChannelID,
			Cause:     "plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	h.logger.Info("Prepare chat channel", "channelID", data.ChannelID, "pluginID", conn.pluginID)
	if _, err := h.sessionMgr.PreparePluginSession(context.Background(), data.ChannelID, conn.pluginID); err != nil {
		h.logger.Error("Prepare chat channel failed", "channelID", data.ChannelID, "err", err)
		_ = conn.Send(contracts.PrepareChatChannelFailedEvent, &contracts.PrepareChatChannelFailedEventPayload{
			ChannelID: data.ChannelID,
			Cause:     err.Error(),
		}, payload.RequestID)
		return
	}

	h.logger.Info("Done prepare chat channel", "channelID", data.ChannelID)
	_ = conn.Send(contracts.PrepareChatChannelSuccessEvent, &contracts.PrepareChatChannelSuccessEventPayload{
		ChannelID: data.ChannelID,
	}, payload.RequestID)
}

func (h *PluginHost) handleSendChatMessage(conn *PluginConnection, payload *contracts.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(contracts.SendChatMessageFailedEvent, &contracts.SendChatMessageFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[contracts.SendChatMessageEventPayload](payload)
	if err != nil {
		_ = conn.Send(contracts.SendChatMessageFailedEvent, &contracts.SendChatMessageFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !RequireScopes(conn.scopes, core_plugin.ChatFeature) {
		_ = conn.Send(contracts.SendChatMessageFailedEvent, &contracts.SendChatMessageFailedEventPayload{
			Cause: "plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	ctx := context.Background()
	session, err := h.sessionMgr.GetPluginSessionByChannel(ctx, data.ChannelID, conn.pluginID)
	if err != nil {
		_ = conn.Send(contracts.SendChatMessageFailedEvent, &contracts.SendChatMessageFailedEventPayload{
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}

	h.eventBus.Publish(ctx, event_bus.SendChatMessageEvent, &event_bus.SendChatMessageEventPayload{
		RequestID: uuid.New().String(),
		SessionID: session.ID.String(),
		ChannelID: data.ChannelID,
		Messages: messages.ListMessage{
			messages.NewTextMessage(messages.USER, data.Text),
		},
		SenderID:   conn.pluginID,
		SenderType: event_bus.Plugin,
	})

	_ = conn.Send(contracts.SendChatMessageSuccessEvent, &contracts.SendChatMessageSuccessEventPayload{
		ChannelID: data.ChannelID,
	}, payload.RequestID)
}

func (h *PluginHost) handleResetChatChannel(conn *PluginConnection, payload *contracts.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(contracts.ResetChatChannelFailedEvent, &contracts.ResetChatChannelFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[contracts.ResetChatChannelEventPayload](payload)
	if err != nil {
		_ = conn.Send(contracts.ResetChatChannelFailedEvent, &contracts.ResetChatChannelFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !RequireScopes(conn.scopes, core_plugin.ChatFeature) {
		_ = conn.Send(contracts.ResetChatChannelFailedEvent, &contracts.ResetChatChannelFailedEventPayload{
			Cause: "plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	err = h.sessionMgr.ResetChatChannel(context.Background(), data.ChannelID, conn.pluginID)
	if err != nil {
		_ = conn.Send(contracts.ResetChatChannelFailedEvent, &contracts.ResetChatChannelFailedEventPayload{
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}

	_ = conn.Send(contracts.ResetChatChannelSuccessEvent, &contracts.ResetChatChannelSuccessEventPayload{
		ChannelID: data.ChannelID,
	}, payload.RequestID)
}

func (h *PluginHost) handleSubmitToolApproval(conn *PluginConnection, payload *contracts.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(contracts.SubmitToolApprovalFailedEvent, &contracts.SubmitToolApprovalFailedEventPayload{Cause: "access denied"}, payload.RequestID)
		return
	}

	data, err := getPayloadData[contracts.SubmitToolApprovalEventPayload](payload)
	if err != nil || data == nil || data.ChannelID == "" || data.ApprovalRequestID == "" {
		_ = conn.Send(contracts.SubmitToolApprovalFailedEvent, &contracts.SubmitToolApprovalFailedEventPayload{Cause: "invalid approval payload"}, payload.RequestID)
		return
	}
	if !RequireScopes(conn.scopes, core_plugin.ChatFeature) {
		_ = conn.Send(contracts.SubmitToolApprovalFailedEvent, &contracts.SubmitToolApprovalFailedEventPayload{Cause: "plugin not supported this feature"}, payload.RequestID)
		return
	}
	session, err := h.sessionMgr.GetPluginSessionByChannel(context.Background(), data.ChannelID, conn.pluginID)
	if err != nil {
		_ = conn.Send(contracts.SubmitToolApprovalFailedEvent, &contracts.SubmitToolApprovalFailedEventPayload{Cause: err.Error()}, payload.RequestID)
		return
	}
	if err := h.approvalMgr.Respond(data.ApprovalRequestID, session.ID.String(), data.Accept); err != nil {
		_ = conn.Send(contracts.SubmitToolApprovalFailedEvent, &contracts.SubmitToolApprovalFailedEventPayload{Cause: err.Error()}, payload.RequestID)
		return
	}

	_ = conn.Send(contracts.SubmitToolApprovalSuccessEvent, &contracts.SubmitToolApprovalSuccessEventPayload{}, payload.RequestID)
}

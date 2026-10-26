package host

import (
	"context"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/event_bus"
	core_plugin "github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/pkgs/messages"
	plugin_dtos "github.com/smtdfc/nagare/plugin/dtos"
)

func (h *PluginHost) handlePrepareChatSession(conn *PluginConnection, payload *plugin_dtos.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin_dtos.PrepareChatSessionFailedEvent, &plugin_dtos.PrepareChatSessionFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin_dtos.PrepareChatSessionEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin_dtos.PrepareChatSessionFailedEvent, &plugin_dtos.PrepareChatSessionFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !RequireScopes(conn.scopes, core_plugin.ChatFeature) {
		_ = conn.Send(plugin_dtos.PrepareChatSessionFailedEvent, &plugin_dtos.PrepareChatSessionFailedEventPayload{
			ChannelID: data.ChannelID,
			Cause:     "plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	h.logger.Info("Prepare chat session", "channelID", data.ChannelID, "pluginID", conn.pluginID)
	session, err := h.sessionMgr.PreparePluginSession(context.Background(), data.ChannelID, conn.pluginID)
	if err != nil {
		h.logger.Error("Prepare chat session failed", "channelID", data.ChannelID, "err", err)
		_ = conn.Send(plugin_dtos.PrepareChatSessionFailedEvent, &plugin_dtos.PrepareChatSessionFailedEventPayload{
			ChannelID: data.ChannelID,
			Cause:     err.Error(),
		}, payload.RequestID)
		return
	}

	h.logger.Info("Done prepare chat session", "channelID", data.ChannelID, "sessionID", session.ID.String())
	_ = conn.Send(plugin_dtos.PrepareChatSessionSuccessEvent, &plugin_dtos.PrepareChatSessionSuccessEventPayload{
		ChannelID: data.ChannelID,
		SessionID: session.ID.String(),
	}, payload.RequestID)
}

func (h *PluginHost) handleSendChatMessage(conn *PluginConnection, payload *plugin_dtos.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin_dtos.SendChatMessageFailedEvent, &plugin_dtos.SendChatMessageFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin_dtos.SendChatMessageEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin_dtos.SendChatMessageFailedEvent, &plugin_dtos.SendChatMessageFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !RequireScopes(conn.scopes, core_plugin.ChatFeature) {
		_ = conn.Send(plugin_dtos.SendChatMessageFailedEvent, &plugin_dtos.SendChatMessageFailedEventPayload{
			Cause: "plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	ctx := context.Background()
	_, err = h.sessionMgr.GetPluginSession(ctx, data.SessionID, conn.pluginID)
	if err != nil {
		_ = conn.Send(plugin_dtos.SendChatMessageFailedEvent, &plugin_dtos.SendChatMessageFailedEventPayload{
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}

	h.chatEventBus.Publish(ctx, event_bus.SendChatMessageEvent, &event_bus.SendChatMessageEventPayload{
		RequestID: uuid.New().String(),
		SessionID: data.SessionID,
		Messages: messages.ListMessage{
			messages.NewTextMessage(messages.USER, data.Text),
		},
		SenderID:   conn.pluginID,
		SenderType: event_bus.Plugin,
	})

	_ = conn.Send(plugin_dtos.SendChatMessageSuccessEvent, &plugin_dtos.SendChatMessageSuccessEventPayload{
		SessionID: data.SessionID,
	}, payload.RequestID)
}

func (h *PluginHost) handleResetChatChannel(conn *PluginConnection, payload *plugin_dtos.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin_dtos.ResetChatChannelFailedEvent, &plugin_dtos.ResetChatChannelFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin_dtos.ResetChatChannelEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin_dtos.ResetChatChannelFailedEvent, &plugin_dtos.ResetChatChannelFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !RequireScopes(conn.scopes, core_plugin.ChatFeature) {
		_ = conn.Send(plugin_dtos.ResetChatChannelFailedEvent, &plugin_dtos.ResetChatChannelFailedEventPayload{
			Cause: "plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	err = h.sessionMgr.ResetChatChannel(context.Background(), data.ChannelID, conn.pluginID)
	if err != nil {
		_ = conn.Send(plugin_dtos.ResetChatChannelFailedEvent, &plugin_dtos.ResetChatChannelFailedEventPayload{
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}

	_ = conn.Send(plugin_dtos.ResetChatChannelSuccessEvent, &plugin_dtos.ResetChatChannelSuccessEventPayload{
		ChannelID: data.ChannelID,
	}, payload.RequestID)
}

func (h *PluginHost) handleSubmitToolApproval(conn *PluginConnection, payload *plugin_dtos.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin_dtos.SubmitToolApprovalFailedEvent, &plugin_dtos.SubmitToolApprovalFailedEventPayload{Cause: "access denied"}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin_dtos.SubmitToolApprovalEventPayload](payload)
	if err != nil || data == nil || data.SessionID == "" || data.ApprovalRequestID == "" {
		_ = conn.Send(plugin_dtos.SubmitToolApprovalFailedEvent, &plugin_dtos.SubmitToolApprovalFailedEventPayload{Cause: "invalid approval payload"}, payload.RequestID)
		return
	}
	if !RequireScopes(conn.scopes, core_plugin.ChatFeature) {
		_ = conn.Send(plugin_dtos.SubmitToolApprovalFailedEvent, &plugin_dtos.SubmitToolApprovalFailedEventPayload{Cause: "plugin not supported this feature"}, payload.RequestID)
		return
	}
	if _, err := h.sessionMgr.GetPluginSession(context.Background(), data.SessionID, conn.pluginID); err != nil {
		_ = conn.Send(plugin_dtos.SubmitToolApprovalFailedEvent, &plugin_dtos.SubmitToolApprovalFailedEventPayload{Cause: err.Error()}, payload.RequestID)
		return
	}
	if err := h.approvalMgr.Respond(data.ApprovalRequestID, data.SessionID, data.Accept); err != nil {
		_ = conn.Send(plugin_dtos.SubmitToolApprovalFailedEvent, &plugin_dtos.SubmitToolApprovalFailedEventPayload{Cause: err.Error()}, payload.RequestID)
		return
	}

	_ = conn.Send(plugin_dtos.SubmitToolApprovalSuccessEvent, &plugin_dtos.SubmitToolApprovalSuccessEventPayload{}, payload.RequestID)
}

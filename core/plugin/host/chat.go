package host

import (
	"context"
	"slices"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/event_bus"
	core_plugin "github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
)

func (h *PluginHost) handlePrepareChatSession(conn *PluginConnection, payload *websocket.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin.PrepareChatSessionFailedEvent, &plugin.PrepareChatSessionFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin.PrepareChatSessionEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin.PrepareChatSessionFailedEvent, &plugin.PrepareChatSessionFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !slices.Contains(conn.scopes, core_plugin.ChatFeature.ToString()) {
		_ = conn.Send(plugin.PrepareChatSessionFailedEvent, &plugin.PrepareChatSessionFailedEventPayload{
			ChannelID: data.ChannelID,
			Cause:     "Plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	h.logger.Info("Prepare chat session", "channelID", data.ChannelID, "pluginID", conn.pluginID)
	session, err := h.sessionMgr.PreparePluginSession(context.Background(), data.ChannelID, conn.pluginID)
	if err != nil {
		h.logger.Error("Prepare chat session failed", "channelID", data.ChannelID, "err", err)
		_ = conn.Send(plugin.PrepareChatSessionFailedEvent, &plugin.PrepareChatSessionFailedEventPayload{
			ChannelID: data.ChannelID,
			Cause:     err.Error(),
		}, payload.RequestID)
		return
	}

	h.logger.Info("Done prepare chat session", "channelID", data.ChannelID, "sessionID", session.ID.String())
	_ = conn.Send(plugin.PrepareChatSessionSuccessEvent, &plugin.PrepareChatSessionSuccessEventPayload{
		ChannelID: data.ChannelID,
		SessionID: session.ID.String(),
	}, payload.RequestID)
}

func (h *PluginHost) handleSendChatMessage(conn *PluginConnection, payload *websocket.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin.SendChatMessageFailedEvent, &plugin.SendChatMessageFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin.SendChatMessageEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin.SendChatMessageFailedEvent, &plugin.SendChatMessageFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !slices.Contains(conn.scopes, core_plugin.ChatFeature.ToString()) {
		_ = conn.Send(plugin.SendChatMessageFailedEvent, &plugin.SendChatMessageFailedEventPayload{
			Cause: "Plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	ctx := context.Background()
	_, err = h.sessionMgr.GetPluginSession(ctx, data.SessionID, conn.pluginID)
	if err != nil {
		_ = conn.Send(plugin.SendChatMessageFailedEvent, &plugin.SendChatMessageFailedEventPayload{
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}

	h.chatEventBus.Publish(ctx, event_bus.SendEvent, &event_bus.SendMessageEventPayload{
		RequestID:  uuid.New().String(),
		SessionID:  data.SessionID,
		Text:       data.Text,
		SenderID:   conn.pluginID,
		SenderType: event_bus.Plugin,
	})

	_ = conn.Send(plugin.SendChatMessageSuccessEvent, &plugin.SendChatMessageSuccessEventPayload{
		SessionID: data.SessionID,
	}, payload.RequestID)
}

func (h *PluginHost) handleResetChatChannel(conn *PluginConnection, payload *websocket.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin.ResetChatChannelFailedEvent, &plugin.ResetChatChannelFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin.ResetChatChannelEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin.ResetChatChannelFailedEvent, &plugin.ResetChatChannelFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !slices.Contains(conn.scopes, core_plugin.ChatFeature.ToString()) {
		_ = conn.Send(plugin.ResetChatChannelFailedEvent, &plugin.ResetChatChannelFailedEventPayload{
			Cause: "Plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	err = h.sessionMgr.ResetChatChannel(context.Background(), data.ChannelID, conn.pluginID)
	if err != nil {
		_ = conn.Send(plugin.ResetChatChannelFailedEvent, &plugin.ResetChatChannelFailedEventPayload{
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}

	_ = conn.Send(plugin.ResetChatChannelSuccessEvent, &plugin.ResetChatChannelSuccessEventPayload{
		ChannelID: data.ChannelID,
	}, payload.RequestID)
}

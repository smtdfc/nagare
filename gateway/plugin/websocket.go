package plugin

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	core_plugin "github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/plugin/manager"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"
	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
	websocket2 "github.com/smtdfc/nagare/gateway/common/websocket"
)

type ChatWebsocketHandler struct {
	pluginMgr         *manager.PluginManager
	sessionMgr        *session_mgr.SessionManager
	chatEventBus      *event_bus.CoreEventBus
	logger            *logger.BaseLogger
	pluginConnections map[string]*melody.Session
}

func (w *ChatWebsocketHandler) OnHandshakeEvent(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
	ctx := context.Background()
	data, err := websocket2.GetData[plugin_dtos.HandshakeEventPayload](message)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			ID:    "",
			Cause: "failed to parsing payload",
		}, message.RequestID)
		return
	}

	w.logger.Info("Handshaking with plugin", "pluginID", data.PluginID)
	pluginInfo, err := w.pluginMgr.ValidConnect(ctx, data.PluginID, data.ConnectCode)
	if err != nil {
		w.logger.Error("Handshake failed", "pluginID", data.PluginID, "cause", err)
		_ = websocket2.SendMessage(s, plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			ID:    data.ID,
			Cause: err.Error(),
		}, message.RequestID)
		return
	}

	scopes := make([]string, 0, len(pluginInfo.Features))
	for _, feature := range pluginInfo.Features {
		scopes = append(scopes, feature.ToString())
	}

	s.Set("auth", &websocket2.AuthData{
		TargetType: "plugin",
		TargetID:   data.PluginID,
		Scopes:     scopes,
	})

	ws.JoinRoom(fmt.Sprintf("plugin:%s:chat", pluginInfo.ID.String()), s)
	err = websocket2.SendMessage(s, plugin_dtos.HandshakeSuccessEvent, &plugin_dtos.HandshakeSuccessEventPayload{
		ID: data.ID,
	}, message.RequestID)
	if err != nil {
		return
	}
	w.logger.Info("Handshake success", "pluginID", data.PluginID)
}

func (w *ChatWebsocketHandler) OnPrepareChatSession(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
	ctx := context.Background()
	data, err := websocket2.GetData[plugin_dtos.PrepareChatSessionEventPayload](message)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.PrepareChatSessionFailedEvent, &plugin_dtos.PrepareChatSessionFailedEventPayload{
			Cause: "failed to parsing payload",
		}, message.RequestID)
		return
	}

	w.logger.Info("debug", "data", data)

	auth, err := w.getPluginAuth(s)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.PrepareChatSessionFailedEvent, &plugin_dtos.PrepareChatSessionFailedEventPayload{
			ChannelID: data.ChannelID,
			Cause:     err.Error(),
		}, message.RequestID)
		return
	}

	w.logger.Info("debug", "data", data)

	if !slices.Contains(auth.Scopes, core_plugin.ChatFeature.ToString()) {
		_ = websocket2.SendMessage(s, plugin_dtos.PrepareChatSessionFailedEvent, &plugin_dtos.PrepareChatSessionFailedEventPayload{
			ChannelID: data.ChannelID,
			Cause:     "Plugin not supported this feature",
		}, message.RequestID)
		return
	}

	w.logger.Info("Prepare chat session", "channelID", data.ChannelID)
	session, err := w.sessionMgr.PreparePluginSession(ctx, data.ChannelID, auth.TargetID)

	if err != nil {
		w.logger.Error("Prepare chat session failed", "channelID", data.ChannelID, "err", err)
		_ = websocket2.SendMessage(s, plugin_dtos.PrepareChatSessionFailedEvent, &plugin_dtos.PrepareChatSessionFailedEventPayload{
			ChannelID: data.ChannelID,
			Cause:     err.Error(),
		}, message.RequestID)
		return
	}

	w.logger.Info("Done ", "channelID", data.ChannelID, "sessionID", session.ID.String())
	_ = websocket2.SendMessage(s, plugin_dtos.PrepareChatSessionSuccessEvent, &plugin_dtos.PrepareChatSessionSuccessEventPayload{
		ChannelID: data.ChannelID,
		SessionID: session.ID.String(),
	}, message.RequestID)

}

func (w *ChatWebsocketHandler) OnSendChatMessage(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
	ctx := context.Background()
	data, err := websocket2.GetData[plugin_dtos.SendChatMessageEventPayload](message)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.SendChatMessageFailedEvent, &plugin_dtos.SendChatMessageFailedEventPayload{
			Cause: "failed to parsing payload",
		}, message.RequestID)
		return
	}
	auth, err := w.getPluginAuth(s)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.SendChatMessageFailedEvent, &plugin_dtos.SendChatMessageFailedEventPayload{
			Cause: err.Error(),
		}, message.RequestID)
		return
	}
	if !slices.Contains(auth.Scopes, core_plugin.ChatFeature.ToString()) {
		_ = websocket2.SendMessage(s, plugin_dtos.SendChatMessageFailedEvent, &plugin_dtos.SendChatMessageFailedEventPayload{
			Cause: "Plugin not supported this feature",
		}, message.RequestID)
		return
	}

	_, err = w.sessionMgr.GetPluginSession(ctx, data.SessionID, auth.TargetID)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.SendChatMessageFailedEvent, &plugin_dtos.SendChatMessageFailedEventPayload{
			Cause: err.Error(),
		}, message.RequestID)
		return
	}

	w.chatEventBus.Publish(ctx, event_bus.SendEvent, &event_bus.SendMessageEventPayload{
		RequestID:  uuid.New().String(),
		SessionID:  data.SessionID,
		Text:       data.Text,
		SenderID:   auth.TargetID,
		SenderType: event_bus.Plugin,
	})

	_ = websocket2.SendMessage(s, plugin_dtos.SendChatMessageSuccessEvent, &plugin_dtos.SendChatMessageSuccessEventPayload{
		SessionID: data.SessionID,
	}, message.RequestID)
}

func (w *ChatWebsocketHandler) OnResetChatChannel(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
	ctx := context.Background()
	data, err := websocket2.GetData[plugin_dtos.ResetChatChannelEventPayload](message)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.ResetChatChannelFailedEvent, &plugin_dtos.ResetChatChannelFailedEventPayload{
			Cause: "failed to parsing payload",
		}, message.RequestID)
		return
	}
	auth, err := w.getPluginAuth(s)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.ResetChatChannelFailedEvent, &plugin_dtos.ResetChatChannelFailedEventPayload{
			Cause: err.Error(),
		}, message.RequestID)
		return
	}

	err = w.sessionMgr.ResetChatChannel(ctx, data.ChannelID, auth.TargetID)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.ResetChatChannelFailedEvent, &plugin_dtos.ResetChatChannelFailedEventPayload{
			Cause: err.Error(),
		}, message.RequestID)
		return
	}

	_ = websocket2.SendMessage(s, plugin_dtos.ResetChatChannelSuccessEvent, &plugin_dtos.ResetChatChannelSuccessEventPayload{
		ChannelID: data.ChannelID,
	}, message.RequestID)
}

func (w *ChatWebsocketHandler) getPluginAuth(s *melody.Session) (*websocket2.AuthData, error) {
	value, exist := s.Get("auth")
	if !exist {
		return nil, errors.New("access denied")
	}

	auth, ok := value.(*websocket2.AuthData)
	if !ok || auth.TargetType != "plugin" {
		return nil, errors.New("access denied")
	}

	return auth, nil
}

// @Injectable
func NewWebsocketHandler(
	pluginMgr *manager.PluginManager,
	sessionMgr *session_mgr.SessionManager,
	chatEventBus *event_bus.CoreEventBus,
	logger *logger.BaseLogger,
) *ChatWebsocketHandler {
	return &ChatWebsocketHandler{
		pluginMgr:    pluginMgr,
		sessionMgr:   sessionMgr,
		chatEventBus: chatEventBus,
		logger:       logger,
	}
}

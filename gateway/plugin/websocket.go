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
	manager2 "github.com/smtdfc/nagare/core/tool/manager"
	tool_registry "github.com/smtdfc/nagare/core/tool/registry"
	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
	websocket2 "github.com/smtdfc/nagare/gateway/common/websocket"
)

type PluginWebsocketHandler struct {
	pluginMgr         *manager.PluginManager
	toolMgr           *manager2.ToolManager
	sessionMgr        *session_mgr.SessionManager
	chatEventBus      *event_bus.CoreEventBus
	logger            *logger.BaseLogger
	pluginConnections map[string]*melody.Session
}

func (w *PluginWebsocketHandler) OnHandshakeEvent(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
	ctx := context.Background()
	data, err := websocket2.GetData[plugin_dtos.HandshakeEventPayload](message)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			ID:    "",
			Cause: "failed to parsing payload",
		}, message.RequestID)
		return
	}

	w.logger.Info("Handshaking with plugin", "packageName", data.PackageName)
	pluginInfo, err := w.pluginMgr.ValidConnect(ctx, data.PackageName, data.ConnectCode)
	if err != nil {
		w.logger.Error("Handshake failed", "packageName", data.PackageName, "cause", err)
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
		TargetID:   pluginInfo.ID.String(),
		Scopes:     scopes,
	})

	ws.JoinRoom(fmt.Sprintf("plugin:%s:chat", pluginInfo.ID.String()), s)
	err = websocket2.SendMessage(s, plugin_dtos.HandshakeSuccessEvent, &plugin_dtos.HandshakeSuccessEventPayload{
		ID: data.ID,
	}, message.RequestID)
	if err != nil {
		return
	}
	w.logger.Info("Handshake success", "packageName", data.PackageName)
}

func (w *PluginWebsocketHandler) OnPrepareChatSession(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
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

func (w *PluginWebsocketHandler) OnSendChatMessage(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
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

func (w *PluginWebsocketHandler) OnResetChatChannel(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
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

	if !slices.Contains(auth.Scopes, core_plugin.ChatFeature.ToString()) {
		_ = websocket2.SendMessage(s, plugin_dtos.ResetChatChannelFailedEvent, &plugin_dtos.ResetChatChannelFailedEventPayload{
			Cause: "Plugin not supported this feature",
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

func (w *PluginWebsocketHandler) OnRegisterTool(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
	data, err := websocket2.GetData[plugin_dtos.RegisterPluginToolEventPayload](message)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.RegisterPluginToolFailedEvent, &plugin_dtos.RegisterPluginToolFailedEventPayload{
			Cause: "failed to parsing payload",
		}, message.RequestID)
		return
	}
	auth, err := w.getPluginAuth(s)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.RegisterPluginToolFailedEvent, &plugin_dtos.RegisterPluginToolFailedEventPayload{
			Cause: err.Error(),
		}, message.RequestID)
		return
	}

	if !slices.Contains(auth.Scopes, core_plugin.PluginToolFeature.ToString()) {
		_ = websocket2.SendMessage(s, plugin_dtos.RegisterPluginToolFailedEvent, &plugin_dtos.RegisterPluginToolFailedEventPayload{
			Cause: "Plugin not supported this feature",
		}, message.RequestID)
		return
	}

	w.logger.Info("Registering plugin tool", "name", data.Tool.Name, "pluginID", auth.TargetID)
	tool_registry.RegisterPluginTool(
		data.Tool.Name,
		data.Tool.Args,
		data.Tool.Description,
		data.Tool.Categories,
		auth.TargetID,
		true,
		w.chatEventBus,
	)

	_ = websocket2.SendMessage(s, plugin_dtos.RegisterPluginToolSuccessEvent, &plugin_dtos.RegisterPluginToolSuccessEventPayload{}, message.RequestID)
}

func (w *PluginWebsocketHandler) OnRegisterToolCategories(s *melody.Session, ws *websocket2.Coordinator, message *websocket.Payload[any]) {
	data, err := websocket2.GetData[plugin_dtos.RegisterToolCategoriesEventPayload](message)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.RegisterToolCategoriesFailedEvent, &plugin_dtos.RegisterToolCategoriesFailedEventPayload{
			Cause: "failed to parsing payload",
		}, message.RequestID)
		return
	}
	auth, err := w.getPluginAuth(s)
	if err != nil {
		_ = websocket2.SendMessage(s, plugin_dtos.RegisterToolCategoriesFailedEvent, &plugin_dtos.RegisterToolCategoriesFailedEventPayload{
			Cause: err.Error(),
		}, message.RequestID)
		return
	}

	if !slices.Contains(auth.Scopes, core_plugin.PluginToolFeature.ToString()) {
		_ = websocket2.SendMessage(s, plugin_dtos.RegisterToolCategoriesFailedEvent, &plugin_dtos.RegisterToolCategoriesFailedEventPayload{
			Cause: "Plugin not supported this feature",
		}, message.RequestID)
		return
	}

	w.toolMgr.AddPluginToolCategories(auth.TargetID, data.Categories)
	w.logger.Info("Registered plugin tool categories", "categories", data.Categories, "pluginID", auth.TargetID)
	_ = websocket2.SendMessage(s, plugin_dtos.RegisterToolCategoriesSuccessEvent, &plugin_dtos.RegisterToolCategoriesSuccessEventPayload{}, message.RequestID)
}

func (w *PluginWebsocketHandler) OnPluginToolCallResult(s *melody.Session, _ *websocket2.Coordinator, message *websocket.Payload[any]) {
	auth, err := w.getPluginAuth(s)
	if err != nil {
		return
	}

	data, err := websocket2.GetData[plugin_dtos.PluginToolCallResultEventPayload](message)
	if err != nil {
		return
	}

	w.chatEventBus.Publish(context.Background(), event_bus.PluginToolCallResultEvent, &event_bus.PluginToolCallResultEventPayload{
		RequestID: message.RequestID,
		Result:    data.Result,
		Error:     data.Error,
	})
	w.logger.Debug("Plugin tool result received", "pluginID", auth.TargetID, "requestID", message.RequestID)
}

func (w *PluginWebsocketHandler) ForwardPluginToolCalls(ws *websocket2.Coordinator) {
	channel, unsubscribe := w.chatEventBus.Subscribe(event_bus.PluginToolCallEvent)
	defer unsubscribe()

	for payload := range channel {
		call, ok := payload.(*event_bus.PluginToolCallEventPayload)
		if !ok {
			continue
		}
		roomID := fmt.Sprintf("plugin:%s:chat", call.PluginID)
		err := websocket2.BroadcastToRoom(ws, roomID, plugin_dtos.PluginToolCallEvent, plugin_dtos.PluginToolCallEventPayload{
			Name: call.Name,
			Args: call.Args,
		}, call.RequestID, nil)
		if err != nil {
			w.logger.Error("Failed to forward plugin tool call", "error", err, "roomID", roomID, "pluginID", call.PluginID, "requestID", call.RequestID)
			continue
		}
		w.logger.Info("Forwarded plugin tool call", "roomID", roomID, "pluginID", call.PluginID, "requestID", call.RequestID)
	}
}

func (w *PluginWebsocketHandler) getPluginAuth(s *melody.Session) (*websocket2.AuthData, error) {
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
	toolMgr *manager2.ToolManager,
	chatEventBus *event_bus.CoreEventBus,
	logger *logger.BaseLogger,
) *PluginWebsocketHandler {
	return &PluginWebsocketHandler{
		pluginMgr:    pluginMgr,
		sessionMgr:   sessionMgr,
		toolMgr:      toolMgr,
		chatEventBus: chatEventBus,
		logger:       logger,
	}
}

package host

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	core_plugin "github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/core/session"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"
	manager2 "github.com/smtdfc/nagare/core/tool/manager"
	tool_registry "github.com/smtdfc/nagare/core/tool/registry"
	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
	"github.com/smtdfc/nagare/pkgs/helpers"
	"github.com/smtdfc/nagare/pkgs/ipc"
	"github.com/smtdfc/nagare/pkgs/paths"
)

type PluginConnection struct {
	conn       net.Conn
	mu         sync.Mutex
	pluginID   string
	pluginName string
	scopes     []string
	closed     bool
}

func (pc *PluginConnection) Send(event websocket.Event, data any, requestID string) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.closed {
		return errors.New("connection closed")
	}

	payload := websocket.Payload[any]{
		Event:     event,
		RequestID: requestID,
		Data:      data,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return ipc.WriteMessage(pc.conn, payloadBytes)
}

func (pc *PluginConnection) Close() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.closed {
		return nil
	}
	pc.closed = true
	return pc.conn.Close()
}

func getPayloadData[T any](payload *websocket.Payload[any]) (*T, error) {
	if payload == nil || payload.Data == nil {
		return nil, errors.New("payload or data is nil")
	}

	bytesData, err := json.Marshal(payload.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload data: %w", err)
	}

	var result T
	if err := json.Unmarshal(bytesData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into target type %T: %w", (*T)(nil), err)
	}

	return &result, nil
}

type PluginHost struct {
	pluginMgr    *manager.PluginManager
	toolMgr      *manager2.ToolManager
	sessionMgr   *session_mgr.SessionManager
	chatEventBus *event_bus.CoreEventBus
	logger       *logger.BaseLogger
	socketPath   string

	listener    net.Listener
	mu          sync.RWMutex
	connections map[string]*PluginConnection
	stopChan    chan struct{}
}

func (h *PluginHost) SetSocketPath(path string) {
	h.socketPath = path
}

func (h *PluginHost) GetSocketPath() string {
	return h.socketPath
}

func (h *PluginHost) Start(ctx context.Context) error {
	path := h.socketPath
	if path == "" {
		path = ipc.GetDefaultSocketPath()
		h.socketPath = path
	}

	listener, err := ipc.Listen(path)
	if err != nil {
		h.logger.Error("Failed to listen on plugin socket", "error", err, "path", path)
		return err
	}

	h.listener = listener
	h.stopChan = make(chan struct{})

	go h.acceptLoop()
	go h.forwardPluginToolCalls()
	go h.forwardChatChunks()

	h.logger.Info("Plugin host started", "socketPath", path)
	return nil
}

func (h *PluginHost) Stop(ctx context.Context) error {
	if h.stopChan != nil {
		select {
		case <-h.stopChan:
		default:
			close(h.stopChan)
		}
	}

	if h.listener != nil {
		_ = h.listener.Close()
	}

	h.mu.Lock()
	for id, conn := range h.connections {
		_ = conn.Close()
		delete(h.connections, id)
	}
	h.mu.Unlock()

	h.logger.Info("Plugin host stopped")
	return nil
}

func (h *PluginHost) acceptLoop() {
	for {
		conn, err := h.listener.Accept()
		if err != nil {
			select {
			case <-h.stopChan:
				return
			default:
				h.logger.Debug("Accept error", "error", err)
				return
			}
		}

		go h.handleConnection(conn)
	}
}

func (h *PluginHost) handleConnection(conn net.Conn) {
	pConn := &PluginConnection{
		conn: conn,
	}

	defer func() {
		_ = pConn.Close()
		if pConn.pluginID != "" {
			h.mu.Lock()
			delete(h.connections, pConn.pluginID)
			h.mu.Unlock()
			h.logger.Info("Plugin disconnected", "pluginID", pConn.pluginID, "packageName", pConn.pluginName)
		}
	}()

	for {
		msgBytes, err := ipc.ReadMessage(conn)
		if err != nil {
			break
		}

		var payload websocket.Payload[any]
		if err := json.Unmarshal(msgBytes, &payload); err != nil {
			h.logger.Error("Failed to unmarshal payload", "error", err)
			continue
		}

		h.dispatch(pConn, &payload)
	}
}

func (h *PluginHost) dispatch(pConn *PluginConnection, payload *websocket.Payload[any]) {
	switch payload.Event {
	case plugin_dtos.HandshakeEvent:
		h.handleHandshake(pConn, payload)
	case plugin_dtos.PrepareChatSessionEvent:
		go h.handlePrepareChatSession(pConn, payload)
	case plugin_dtos.SendChatMessageEvent:
		go h.handleSendChatMessage(pConn, payload)
	case plugin_dtos.ResetChatChannelEvent:
		go h.handleResetChatChannel(pConn, payload)
	case plugin_dtos.RegisterPluginToolEvent:
		go h.handleRegisterPluginTool(pConn, payload)
	case plugin_dtos.RegisterToolCategoriesEvent:
		go h.handleRegisterToolCategories(pConn, payload)
	case plugin_dtos.PluginToolCallResultEvent:
		h.handlePluginToolCallResult(pConn, payload)
	default:
		h.logger.Warn("Unknown event received from plugin", "event", payload.Event, "pluginID", pConn.pluginID)
	}
}

func (h *PluginHost) handleHandshake(conn *PluginConnection, payload *websocket.Payload[any]) {
	data, err := getPayloadData[plugin_dtos.HandshakeEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			ID:    "",
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	h.logger.Info("Handshaking with plugin", "packageName", data.PackageName)
	pluginInfo, err := h.pluginMgr.ValidConnect(context.Background(), data.PackageName, data.ConnectCode)
	if err != nil {
		h.logger.Error("Handshake failed", "packageName", data.PackageName, "cause", err)
		_ = conn.Send(plugin_dtos.HandshakeFailedEvent, &plugin_dtos.HandshakeFailedEventPayload{
			ID:    data.ID,
			Cause: err.Error(),
		}, payload.RequestID)
		return
	}

	scopes := make([]string, 0, len(pluginInfo.Features))
	for _, feature := range pluginInfo.Features {
		scopes = append(scopes, feature.ToString())
	}

	conn.pluginID = pluginInfo.ID.String()
	conn.pluginName = pluginInfo.PackageName
	conn.scopes = scopes

	h.mu.Lock()
	h.connections[conn.pluginID] = conn
	h.mu.Unlock()

	err = conn.Send(plugin_dtos.HandshakeSuccessEvent, &plugin_dtos.HandshakeSuccessEventPayload{
		ID: data.ID,
	}, payload.RequestID)
	if err != nil {
		return
	}
	h.logger.Info("Handshake success", "packageName", data.PackageName, "pluginID", conn.pluginID)
}

func (h *PluginHost) handlePrepareChatSession(conn *PluginConnection, payload *websocket.Payload[any]) {
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

	if !slices.Contains(conn.scopes, core_plugin.ChatFeature.ToString()) {
		_ = conn.Send(plugin_dtos.PrepareChatSessionFailedEvent, &plugin_dtos.PrepareChatSessionFailedEventPayload{
			ChannelID: data.ChannelID,
			Cause:     "Plugin not supported this feature",
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

func (h *PluginHost) handleSendChatMessage(conn *PluginConnection, payload *websocket.Payload[any]) {
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

	if !slices.Contains(conn.scopes, core_plugin.ChatFeature.ToString()) {
		_ = conn.Send(plugin_dtos.SendChatMessageFailedEvent, &plugin_dtos.SendChatMessageFailedEventPayload{
			Cause: "Plugin not supported this feature",
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

	h.chatEventBus.Publish(ctx, event_bus.SendEvent, &event_bus.SendMessageEventPayload{
		RequestID:  uuid.New().String(),
		SessionID:  data.SessionID,
		Text:       data.Text,
		SenderID:   conn.pluginID,
		SenderType: event_bus.Plugin,
	})

	_ = conn.Send(plugin_dtos.SendChatMessageSuccessEvent, &plugin_dtos.SendChatMessageSuccessEventPayload{
		SessionID: data.SessionID,
	}, payload.RequestID)
}

func (h *PluginHost) handleResetChatChannel(conn *PluginConnection, payload *websocket.Payload[any]) {
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

	if !slices.Contains(conn.scopes, core_plugin.ChatFeature.ToString()) {
		_ = conn.Send(plugin_dtos.ResetChatChannelFailedEvent, &plugin_dtos.ResetChatChannelFailedEventPayload{
			Cause: "Plugin not supported this feature",
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

func (h *PluginHost) handleRegisterPluginTool(conn *PluginConnection, payload *websocket.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin_dtos.RegisterPluginToolFailedEvent, &plugin_dtos.RegisterPluginToolFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin_dtos.RegisterPluginToolEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin_dtos.RegisterPluginToolFailedEvent, &plugin_dtos.RegisterPluginToolFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !slices.Contains(conn.scopes, core_plugin.PluginToolFeature.ToString()) {
		_ = conn.Send(plugin_dtos.RegisterPluginToolFailedEvent, &plugin_dtos.RegisterPluginToolFailedEventPayload{
			Cause: "Plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	h.logger.Info("Registering plugin tool", "name", data.Tool.Name, "pluginID", conn.pluginID)
	tool_registry.RegisterPluginTool(
		data.Tool.Name,
		data.Tool.Args,
		data.Tool.Description,
		data.Tool.Categories,
		conn.pluginID,
		true,
		h.chatEventBus,
	)

	_ = conn.Send(plugin_dtos.RegisterPluginToolSuccessEvent, &plugin_dtos.RegisterPluginToolSuccessEventPayload{}, payload.RequestID)
}

func (h *PluginHost) handleRegisterToolCategories(conn *PluginConnection, payload *websocket.Payload[any]) {
	if conn.pluginID == "" {
		_ = conn.Send(plugin_dtos.RegisterToolCategoriesFailedEvent, &plugin_dtos.RegisterToolCategoriesFailedEventPayload{
			Cause: "access denied",
		}, payload.RequestID)
		return
	}

	data, err := getPayloadData[plugin_dtos.RegisterToolCategoriesEventPayload](payload)
	if err != nil {
		_ = conn.Send(plugin_dtos.RegisterToolCategoriesFailedEvent, &plugin_dtos.RegisterToolCategoriesFailedEventPayload{
			Cause: "failed to parse payload",
		}, payload.RequestID)
		return
	}

	if !slices.Contains(conn.scopes, core_plugin.PluginToolFeature.ToString()) {
		_ = conn.Send(plugin_dtos.RegisterToolCategoriesFailedEvent, &plugin_dtos.RegisterToolCategoriesFailedEventPayload{
			Cause: "Plugin not supported this feature",
		}, payload.RequestID)
		return
	}

	h.toolMgr.AddPluginToolCategories(conn.pluginID, data.Categories)
	h.logger.Info("Registered plugin tool categories", "categories", data.Categories, "pluginID", conn.pluginID)
	_ = conn.Send(plugin_dtos.RegisterToolCategoriesSuccessEvent, &plugin_dtos.RegisterToolCategoriesSuccessEventPayload{}, payload.RequestID)
}

func (h *PluginHost) handlePluginToolCallResult(conn *PluginConnection, payload *websocket.Payload[any]) {
	data, err := getPayloadData[plugin_dtos.PluginToolCallResultEventPayload](payload)
	if err != nil {
		return
	}

	h.chatEventBus.Publish(context.Background(), event_bus.PluginToolCallResultEvent, &event_bus.PluginToolCallResultEventPayload{
		RequestID: payload.RequestID,
		Result:    data.Result,
		Error:     data.Error,
	})
	h.logger.Debug("Plugin tool result received", "pluginID", conn.pluginID, "requestID", payload.RequestID)
}

func (h *PluginHost) forwardPluginToolCalls() {
	channel, unsubscribe := h.chatEventBus.Subscribe(event_bus.PluginToolCallEvent)
	defer unsubscribe()

	for {
		select {
		case <-h.stopChan:
			return
		case payload, ok := <-channel:
			if !ok {
				return
			}
			call, ok := payload.(*event_bus.PluginToolCallEventPayload)
			if !ok {
				continue
			}

			h.mu.RLock()
			conn, exists := h.connections[call.PluginID]
			h.mu.RUnlock()

			if !exists {
				h.logger.Error("Plugin connection not found for tool call", "pluginID", call.PluginID, "name", call.Name, "requestID", call.RequestID)
				continue
			}

			err := conn.Send(plugin_dtos.PluginToolCallEvent, plugin_dtos.PluginToolCallEventPayload{
				Name: call.Name,
				Args: call.Args,
			}, call.RequestID)
			if err != nil {
				h.logger.Error("Failed to forward plugin tool call", "error", err, "pluginID", call.PluginID, "requestID", call.RequestID)
				continue
			}
			h.logger.Info("Forwarded plugin tool call", "pluginID", call.PluginID, "requestID", call.RequestID)
		}
	}
}

func (h *PluginHost) forwardChatChunks() {
	ch, unsubscribe := h.chatEventBus.Subscribe(event_bus.ChunkEvent)
	defer unsubscribe()

	for {
		select {
		case <-h.stopChan:
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			if evt.GetEventType() != event_bus.ChunkEvent {
				continue
			}
			chunkPayload, ok := evt.(*event_bus.ChatChunkEventPayload)
			if !ok {
				continue
			}

			if chunkPayload.SessionOwnerType != string(session.PLUGIN) && chunkPayload.SessionOwnerType != string(event_bus.Plugin) {
				continue
			}

			h.mu.RLock()
			conn, exists := h.connections[chunkPayload.SessionOwnerID]
			h.mu.RUnlock()

			if !exists {
				continue
			}

			chunkJson, err := helpers.MarshalJson(chunkPayload.Chunk)
			if err != nil {
				continue
			}

			_ = conn.Send(
				websocket.ReceivedChatMessageEvent,
				&websocket.ReceivedChatMessageEventPayload{
					SessionID: chunkPayload.SessionID,
					Message:   chunkJson,
					ChannelID: chunkPayload.ChannelID,
				},
				chunkPayload.RequestID,
			)
		}
	}
}

// @Injectable
func NewPluginHost(
	pluginMgr *manager.PluginManager,
	sessionMgr *session_mgr.SessionManager,
	toolMgr *manager2.ToolManager,
	chatEventBus *event_bus.CoreEventBus,
	logger *logger.BaseLogger,
) *PluginHost {
	return &PluginHost{
		pluginMgr:    pluginMgr,
		sessionMgr:   sessionMgr,
		toolMgr:      toolMgr,
		chatEventBus: chatEventBus,
		logger:       logger.With("module", "plugin-host"),
		connections:  make(map[string]*PluginConnection),
		socketPath:   paths.PluginSocketPath,
	}
}

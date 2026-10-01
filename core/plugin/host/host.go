package host

import (
	"context"
	"encoding/json"
	"net"
	"sync"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/plugin/manager"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"
	manager2 "github.com/smtdfc/nagare/core/tool/manager"
	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
	"github.com/smtdfc/nagare/pkgs/ipc"
	"github.com/smtdfc/nagare/pkgs/paths"
)

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

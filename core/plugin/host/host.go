package host

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/hil"
	"github.com/smtdfc/nagare/core/logger"
	provider_mgr "github.com/smtdfc/nagare/core/model_provider/manager"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/core/prompt"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"
	manager2 "github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/pkgs/ipc"
	"github.com/smtdfc/nagare/pkgs/paths"
	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

type PluginHost struct {
	modelProvider *provider_mgr.ModelProviderManager
	pluginMgr     *manager.PluginManager
	toolMgr       *manager2.ToolManager
	sessionMgr    *session_mgr.SessionManager
	eventBus      *event_bus.CoreEventBus
	approvalMgr   *hil.Manager
	logger        *logger.BaseLogger
	socketPath    string

	listener      net.Listener
	mu            sync.RWMutex
	connections   map[string]*PluginConnection
	chatListeners map[string]*PluginConnection
	stopChan      chan struct{}
}

func (h *PluginHost) unregister(conn *PluginConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.connections[conn.pluginID] == conn {
		delete(h.connections, conn.pluginID)
		prompt.ClearPluginSystemPrompt(conn.pluginID)
		for sessionID, c := range h.chatListeners {
			if c == conn {
				delete(h.chatListeners, sessionID)
			}
		}
	}
}

func (h *PluginHost) SetSocketPath(path string) {
	h.socketPath = path
}

func (h *PluginHost) GetSocketPath() string {
	return h.socketPath
}

func (h *PluginHost) Start() error {
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

func (h *PluginHost) Stop() error {
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
	conns := make([]*PluginConnection, 0, len(h.connections))
	for _, c := range h.connections {
		conns = append(conns, c)
	}
	h.mu.Unlock()

	for _, c := range conns {
		_ = c.Close()
		h.unregister(c)
	}

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
			h.unregister(pConn)
			h.logger.Info("Plugin disconnected", "pluginID", pConn.pluginID, "packageName", pConn.pluginName)
		}
	}()

	for {
		msgBytes, err := ipc.ReadMessage(conn)
		if err != nil {
			break
		}

		var payload contracts.Payload[any]
		if err := json.Unmarshal(msgBytes, &payload); err != nil {
			h.logger.Error("Failed to unmarshal payload", "error", err)
			continue
		}

		h.dispatch(pConn, &payload)
	}
}

func (h *PluginHost) dispatch(pConn *PluginConnection, payload *contracts.Payload[any]) {
	switch payload.Event {
	case contracts.HandshakeEvent:
		h.handleHandshake(pConn, payload)
	case contracts.PrepareChatChannelEvent:
		go h.handlePrepareChatChannel(pConn, payload)
	case contracts.SendChatMessageEvent:
		go h.handleSendChatMessage(pConn, payload)
	case contracts.ResetChatChannelEvent:
		go h.handleResetChatChannel(pConn, payload)
	case contracts.SubmitToolApprovalEvent:
		go h.handleSubmitToolApproval(pConn, payload)
	case contracts.RegisterPluginToolEvent:
		go h.handleRegisterPluginTool(pConn, payload)
	case contracts.RegisterToolCategoriesEvent:
		go h.handleRegisterToolCategories(pConn, payload)
	case contracts.PluginToolCallResultEvent:
		h.handlePluginToolCallResult(pConn, payload)
	case contracts.GetModelProvidersEvent:
		h.handleGetModelProviders(pConn, payload)
	case contracts.RegisterChatListenerEvent:
		h.handleRegisterChatListener(pConn, payload)
	case contracts.InjectSystemPromptEvent:
		go h.handleInjectSystemPrompt(pConn, payload)
	case contracts.ClearSystemPromptEvent:
		go h.handleClearSystemPrompt(pConn, payload)
	default:
		h.logger.Warn("Unknown event received from plugin", "event", payload.Event, "pluginID", pConn.pluginID)
	}
}

// @Injectable
func NewPluginHost(
	pluginMgr *manager.PluginManager,
	sessionMgr *session_mgr.SessionManager,
	toolMgr *manager2.ToolManager,
	modelProvider *provider_mgr.ModelProviderManager,
	eventBus *event_bus.CoreEventBus,
	approvalMgr *hil.Manager,
	logger *logger.BaseLogger,
) *PluginHost {
	return &PluginHost{
		pluginMgr:     pluginMgr,
		sessionMgr:    sessionMgr,
		toolMgr:       toolMgr,
		eventBus:      eventBus,
		approvalMgr:   approvalMgr,
		modelProvider: modelProvider,
		logger:        logger.With("module", "plugin-host"),
		connections:   make(map[string]*PluginConnection),
		chatListeners: make(map[string]*PluginConnection),
		socketPath:    paths.PluginSocketPath,
	}
}

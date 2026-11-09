package external_module

import (
	"encoding/json"
	"net"
	"strings"
	"sync"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/pkgs/ipc"
	"github.com/smtdfc/nagare/pkgs/paths"
)

type ExternalModuleHost struct {
	eventBus          *event_bus.CoreEventBus
	externalModuleMgr *ExternalModuleManager
	logger            *logger.BaseLogger
	socketPath        string

	listener    net.Listener
	listeners   map[string]net.Listener
	mu          sync.RWMutex
	connections map[string]*Connection
	stopChan    chan struct{}

	invokeChannels map[string]chan InvokePayload[any]
	replyChannels  map[string]chan ReplyPayload[any]
}

func (h *ExternalModuleHost) unregister(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.connections[conn.moduleID] == conn {
		delete(h.connections, conn.moduleID)
	}
}

func (h *ExternalModuleHost) SetSocketPath(path string) {
	h.socketPath = path
}

func (h *ExternalModuleHost) GetSocketPath() string {
	return h.socketPath
}

func (h *ExternalModuleHost) Start() error {
	path := h.socketPath
	if path == "" {
		path = ipc.GetDefaultSocketPath()
		h.socketPath = path
	}

	listener, err := ipc.Listen(path)
	if err != nil {
		h.logger.Error("Failed to listen on external module default socket", "error", err, "path", path)
		return err
	}

	h.listener = listener
	h.stopChan = make(chan struct{})

	go h.acceptLoop(listener, "")

	h.logger.Info("External module host started", "defaultSocketPath", path)
	return nil
}

// StartModuleListener tạo một socket/pipe riêng độc quyền cho một moduleID cụ thể
func (h *ExternalModuleHost) StartModuleListener(moduleID string) error {
	socketPath := paths.GetModuleSocketPath(moduleID)

	listener, err := ipc.Listen(socketPath)
	if err != nil {
		h.logger.Error("Failed to listen on dedicated module socket", "moduleID", moduleID, "error", err, "path", socketPath)
		return err
	}

	h.mu.Lock()
	if h.listeners == nil {
		h.listeners = make(map[string]net.Listener)
	}
	h.listeners[moduleID] = listener
	h.mu.Unlock()

	go h.acceptLoop(listener, moduleID)

	h.logger.Info("Dedicated module socket listening", "moduleID", moduleID, "path", socketPath)
	return nil
}

func (h *ExternalModuleHost) Stop() error {
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
	for modID, l := range h.listeners {
		_ = l.Close()
		delete(h.listeners, modID)
	}

	conns := make([]*Connection, 0, len(h.connections))
	for _, c := range h.connections {
		conns = append(conns, c)
	}
	h.mu.Unlock()

	for _, c := range conns {
		_ = c.Close()
		h.unregister(c)
	}

	h.logger.Info("External module host stopped")
	return nil
}

func (h *ExternalModuleHost) acceptLoop(listener net.Listener, expectedModuleID string) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-h.stopChan:
				return
			default:
				h.logger.Debug("Accept error", "error", err, "expectedModuleID", expectedModuleID)
				return
			}
		}

		go h.handleConnection(conn, expectedModuleID)
	}
}

func (h *ExternalModuleHost) handleConnection(conn net.Conn, expectedModuleID string) {
	pConn := &Connection{
		conn:     conn,
		moduleID: expectedModuleID,
	}

	defer func() {
		_ = pConn.Close()
		if pConn.moduleID != "" {
			h.unregister(pConn)
			h.logger.Info("External module disconnected", "moduleID", pConn.moduleID)
		}
	}()

	for {
		msgBytes, err := ipc.ReadMessage(conn)
		if err != nil {
			break
		}

		var payload Payload[any]
		if err := json.Unmarshal(msgBytes, &payload); err != nil {
			h.logger.Error("Failed to unmarshal payload", "error", err)
			continue
		}

		h.dispatch(pConn, &payload)
	}
}

func (h *ExternalModuleHost) dispatch(pConn *Connection, payload *Payload[any]) {
	switch {
	case payload.Event == HandshakeEvent:
		h.onHandshake(pConn, payload)
	case strings.HasPrefix(string(payload.Event), ReplyPrefix):
		h.forwardToInternal(pConn, payload)
	default:
		h.logger.Warn("Unknown event received from external module", "event", payload.Event, "moduleID", pConn.moduleID)
	}
}

func (h *ExternalModuleHost) StartAllModules() error {
	_ = h.Start()

	for moduleID, module := range h.externalModuleMgr.ExternalModules {
		if err := h.StartModuleListener(moduleID); err != nil {
			h.logger.Error("Failed to start dedicated listener for module", "moduleID", moduleID, "error", err)
			continue
		}

		h.logger.Info("Starting external module", "external-module", module.ID)
		connectCode := h.externalModuleMgr.CreateConnectCode(moduleID)
		err := module.Start(connectCode)
		if err != nil {
			h.logger.Error("Failed to start external module", "external-module", module.ID, "error", err)
		}
	}

	return nil
}

func (h *ExternalModuleHost) StopAllModules() error {
	for _, module := range h.externalModuleMgr.ExternalModules {
		h.logger.Info("Stopping external module", "external-module", module.ID)
		err := module.Stop()
		if err != nil {
			h.logger.Error("Failed to stop external module", "external-module", module.ID)
		}
	}

	err := h.Stop()
	if err != nil {
		return err
	}

	return nil
}

// @Injectable
func NewExternalModuleHost(
	eventBus *event_bus.CoreEventBus,
	logger *logger.BaseLogger,
	externalModuleMgr *ExternalModuleManager,
) *ExternalModuleHost {
	return &ExternalModuleHost{
		eventBus:          eventBus,
		logger:            logger.With("module", "external-module-host"),
		connections:       make(map[string]*Connection),
		listeners:         make(map[string]net.Listener),
		socketPath:        paths.ModuleSocketPath,
		externalModuleMgr: externalModuleMgr,
		invokeChannels:    make(map[string]chan InvokePayload[any]),
		replyChannels:     make(map[string]chan ReplyPayload[any]),
	}
}

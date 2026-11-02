package external_module

import (
	"encoding/json"
	"net"
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
	mu          sync.RWMutex
	connections map[string]*Connection
	stopChan    chan struct{}
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
		h.logger.Error("Failed to listen on external module socket", "error", err, "path", path)
		return err
	}

	h.listener = listener
	h.stopChan = make(chan struct{})

	go h.acceptLoop()

	h.logger.Info("External module host started", "socketPath", path)
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

func (h *ExternalModuleHost) acceptLoop() {
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

func (h *ExternalModuleHost) handleConnection(conn net.Conn) {
	pConn := &Connection{
		conn: conn,
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
	switch payload.Event {
	case HandshakeEvent:
		h.onHandshake(pConn, payload)
	default:
		h.logger.Warn("Unknown event received from external module", "event", payload.Event, "moduleID", pConn.moduleID)
	}
}

func (h *ExternalModuleHost) StartAllModules() error {
	err := h.Start()
	if err != nil {
		return err
	}

	for moduleID, module := range h.externalModuleMgr.ExternalModules {
		h.logger.Info("Starting external module", "external-module", module.ID)
		connectCode := h.externalModuleMgr.CreateConnectCode(moduleID)
		err := module.Start(connectCode)
		if err != nil {
			h.logger.Error("Failed to start external module", "external-module", module.ID)
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
		socketPath:        paths.ModuleSocketPath,
		externalModuleMgr: externalModuleMgr,
	}
}

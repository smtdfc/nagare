package host

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"

	"github.com/smtdfc/nagare/plugin-sdk/shared"
)

func GetFreePort(startPort int) (int, error) {
	for port := startPort; port <= 65535; port++ {
		address := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", address)
		if err == nil {
			ln.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("No port free")
}

type Handler func(shared.Message)
type Host struct {
	port     int
	ln       net.Listener
	plugins  map[string]net.Conn
	mu       sync.RWMutex
	logger   *slog.Logger
	handlers map[string]Handler
}

func (h *Host) Handler(kind string, handler Handler) {
	h.mu.Lock()
	h.handlers[kind] = handler
	h.mu.Unlock()
}

func (h *Host) Register(id string, conn net.Conn) {
	h.mu.Lock()
	h.plugins[id] = conn
	h.mu.Unlock()
}

func (h *Host) Send(id string, method string, payload interface{}) error {
	h.mu.RLock()
	conn := h.plugins[id]
	h.mu.RUnlock()

	if conn != nil {
		data, _ := json.Marshal(payload)
		msg := shared.Message{
			PluginID: id,
			Kind:     method,
			Payload:  data,
		}
		json.NewEncoder(conn).Encode(msg)
	}

	return nil
}

func handlePlugin(h *Host, conn net.Conn) {
	var pluginID string
	defer func() {
		conn.Close()
		if pluginID != "" {
			h.mu.Lock()
			delete(h.plugins, pluginID)
			h.mu.Unlock()
		}
	}()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var msg shared.Message
		if err := decoder.Decode(&msg); err != nil {
			return
		}

		switch msg.Kind {
		case shared.REGISTER_PLUGIN_REQUEST:
			json.Unmarshal(msg.Payload, &pluginID)
			h.Register(pluginID, conn)
			encoder.Encode(shared.Message{
				Kind: shared.REGISTER_PLUGIN_SUCCESS,
			})
			h.logger.Info("Plugin registered", "plugin_id", pluginID)
			continue

		case "log":
			encoder.Encode(shared.Message{
				Kind: "plugin:log:success",
			})
			continue
		}

		handler, ok := h.handlers[msg.Kind]
		if ok {
			handler(msg)
		}
	}
}

func (h *Host) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, conn := range h.plugins {
		encoder := json.NewEncoder(conn)
		_ = encoder.Encode(shared.Message{
			Kind: shared.SHUTDOWN_PLUGIN_REQUEST,
		})

		conn.Close()
	}

	h.plugins = nil
	h.logger.Info("Plugin host Shutdown")
}

func (h *Host) Start() {
	ln, _ := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", h.port))
	h.ln = ln
	h.logger.Info("Plugin host started")
	for {
		conn, _ := h.ln.Accept()
		h.logger.Info("New connection")
		go handlePlugin(h, conn)
	}
}

func CreateConnectionString(
	hostname string,
	port int,
) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d", hostname, port)))
	return encoded
}

func NewHost(logger *slog.Logger) *Host {
	port, _ := GetFreePort(3000)
	os.Setenv(shared.ADDR_ENV, CreateConnectionString("localhost", port))

	return &Host{
		port:     port,
		ln:       nil,
		plugins:  map[string]net.Conn{},
		handlers: map[string]Handler{},
		mu:       sync.RWMutex{},
		logger:   logger,
	}
}

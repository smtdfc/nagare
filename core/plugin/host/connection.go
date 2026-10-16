package host

import (
	"encoding/json"
	"errors"
	"net"
	"sync"

	"github.com/smtdfc/nagare/pkgs/ipc"
	plugin_dtos "github.com/smtdfc/nagare/plugin/dtos"
)

type PluginConnection struct {
	conn       net.Conn
	mu         sync.Mutex
	pluginID   string
	pluginName string
	scopes     []string
	closed     bool
}

func (pc *PluginConnection) Send(event plugin_dtos.Event, data any, requestID string) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.closed {
		return errors.New("connection closed")
	}

	payload := plugin_dtos.Payload[any]{
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

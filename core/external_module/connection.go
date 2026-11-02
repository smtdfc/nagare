package external_module

import (
	"encoding/json"
	"errors"
	"net"
	"sync"

	"github.com/smtdfc/nagare/pkgs/ipc"
)

type Connection struct {
	conn       net.Conn
	mu         sync.Mutex
	moduleID   string
	moduleName string
	scopes     []string
	closed     bool
}

func (pc *Connection) Send(event Event, data any, requestID string) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.closed {
		return errors.New("connection closed")
	}

	payload := Payload[any]{
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

func (pc *Connection) Close() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.closed {
		return nil
	}
	pc.closed = true
	return pc.conn.Close()
}

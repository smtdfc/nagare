package client

import (
	"context"
	"encoding/json"
	"net"
	"sync"

	"github.com/smtdfc/nagare/dtos/websocket"
	"github.com/smtdfc/nagare/pkgs/ipc"
)

type Connector struct {
	conn     net.Conn
	mu       sync.Mutex
	onEvent  func(*websocket.Payload[any])
	done     chan struct{}
	doneOnce sync.Once
}

func NewConnector(onEvent func(*websocket.Payload[any])) *Connector {
	return &Connector{
		onEvent: onEvent,
		done:    make(chan struct{}),
	}
}

func (c *Connector) Connect(ctx context.Context, path string) error {
	conn, err := ipc.Dial(ctx, path)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.listenLoop()

	return nil
}

func (c *Connector) listenLoop() {
	defer func() {
		_ = c.Close()
		c.doneOnce.Do(func() { close(c.done) })
	}()

	for {
		msg, err := ipc.ReadMessage(c.conn)
		if err != nil {
			break
		}

		if c.onEvent != nil {
			var payload websocket.Payload[any]
			if err := json.Unmarshal(msg, &payload); err != nil {
				continue
			}
			c.onEvent(&payload)
		}
	}
}

func (c *Connector) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.done:
		return nil
	}
}

func (c *Connector) Send(eventType websocket.Event, data any, requestID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return ErrConnectionNotReady
	}

	payload := websocket.Payload[any]{
		Event:     eventType,
		Data:      data,
		RequestID: requestID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return ErrCreatePayloadFailed
	}

	return ipc.WriteMessage(c.conn, payloadBytes)
}

func (c *Connector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

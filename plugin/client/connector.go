package client

import (
	"context"
	"encoding/json"
	"net"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/smtdfc/nagare/dtos/websocket"
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

func (c *Connector) Connect(ctx context.Context, url string) error {
	conn, _, _, err := ws.DefaultDialer.Dial(ctx, url)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.listenLoop()

	return nil
}

func (c *Connector) listenLoop() {
	defer func() {
		c.Close()
		c.doneOnce.Do(func() { close(c.done) })
	}()

	for {
		msg, _, err := wsutil.ReadServerData(c.conn)
		if err != nil {
			break
		}

		if c.onEvent != nil {
			var payload websocket.Payload[any]
			if err := json.Unmarshal(msg, &payload); err != nil {
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

	return wsutil.WriteClientText(c.conn, payloadBytes)
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

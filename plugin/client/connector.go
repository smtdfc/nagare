package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/smtdfc/nagare/shared/dtos/websocket"
)

func GetData[T any](payload *websocket.Payload[any]) (*T, error) {
	if payload == nil || payload.Data == nil {
		return nil, fmt.Errorf("payload or data is nil")
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

type Connector struct {
	conn    net.Conn
	mu      sync.Mutex
	onEvent func(*websocket.Payload[any])
}

func NewConnector(onEvent func(*websocket.Payload[any])) *Connector {
	return &Connector{
		onEvent: onEvent,
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
	defer c.Close()

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

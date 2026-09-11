package client

import (
	"context"
	"time"
	"uuid"

	plugin_dtos "github.com/smtdfc/nagare/shared/dtos/plugin"
	"github.com/smtdfc/nagare/shared/dtos/websocket"
)

func (p *PluginClient) Handshake(ctx context.Context) error {
	requestID := uuid.NewV4().String()
	respChan := make(chan *websocket.Payload[any], 1)

	p.mu.Lock()
	p.pendingRequests[requestID] = respChan
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		delete(p.pendingRequests, requestID)
		p.mu.Unlock()
	}()

	err := p.connector.Send(
		plugin_dtos.HandshakeEvent,
		plugin_dtos.HandshakeEventPayload{
			ID:          requestID,
			PluginID:    p.Metadata.ID,
			ConnectCode: p.Config.ConnectCode,
		},
		requestID,
	)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ErrHandshakeCancelled
	case <-time.After(5 * time.Second):
		return ErrHandshakeTimeout
	case resp := <-respChan:
		if resp.Event == plugin_dtos.HandshakeFailedEvent {
			return ErrHandshakeFailed
		}
		return nil
	}
}

func (p *PluginClient) onHandshakeSuccess(raw *websocket.Payload[any]) {
	p.mu.Lock()
	if ch, exists := p.pendingRequests[raw.RequestID]; exists {
		ch <- raw
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()
}

func (p *PluginClient) onHandshakeFailed(raw *websocket.Payload[any]) {
	p.mu.Lock()
	if ch, exists := p.pendingRequests[raw.RequestID]; exists {
		ch <- raw
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()
}

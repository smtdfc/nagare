package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/smtdfc/nagare/plugin/metadata"
	plugin_dtos "github.com/smtdfc/nagare/shared/dtos/plugin"
	"github.com/smtdfc/nagare/shared/dtos/websocket"
)

type PluginClient struct {
	Metadata        *metadata.PluginMetadata
	Config          *Config
	connector       *Connector
	mu              sync.Mutex
	pendingRequests map[string]chan *websocket.Payload[any]
}

func (p *PluginClient) Start(ctx context.Context, onStart func()) error {
	p.Config.Port = os.Getenv("NAGARE_PLUGIN_HOST_PORT")
	p.Config.ConnectCode = os.Getenv("NAGARE_PLUGIN_CONNECT_CODE")

	err := p.connector.Connect(ctx, fmt.Sprintf("ws://127.0.0.1:%s/ws", p.Config.Port))
	if err != nil {
		return err
	}

	onStart()
	return nil
}

func (p *PluginClient) handleEvent(payload *websocket.Payload[any]) {
	switch payload.Event {
	case plugin_dtos.HandshakeSuccessEvent:
		p.onHandshakeSuccess(payload)
	case plugin_dtos.HandshakeFailedEvent:
		p.onHandshakeFailed(payload)
	default:
	}
}

func NewPlugin() *PluginClient {
	p := &PluginClient{
		Config:          &Config{},
		pendingRequests: make(map[string]chan *websocket.Payload[any]),
	}

	p.connector = NewConnector(p.handleEvent)
	return p
}

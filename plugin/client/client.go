package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
	"github.com/smtdfc/nagare/plugin/metadata"
	"gopkg.in/natefinch/lumberjack.v2"
)

type PluginClient struct {
	Metadata              *metadata.PluginMetadata
	Config                *Config
	Logger                *slog.Logger
	connector             *Connector
	mu                    sync.Mutex
	pendingRequests       map[string]chan *websocket.Payload[any]
	OnReceivedChatMessage func(sessionID string, channelID string, chunk string)
}

func (p *PluginClient) Start(ctx context.Context, onStart func()) error {
	p.Config.Port = os.Getenv("NAGARE_PLUGIN_HOST_PORT")
	p.Config.ConnectCode = os.Getenv("NAGARE_PLUGIN_CONNECT_CODE")

	err := p.connector.Connect(ctx, fmt.Sprintf("ws://127.0.0.1:%s/ws", p.Config.Port))
	if err != nil {
		p.Logger.Error("Start plugin error", "error", err)
		return err
	}

	onStart()
	return nil
}

func (p *PluginClient) handleEvent(payload *websocket.Payload[any]) {
	switch payload.Event {
	case plugin_dtos.HandshakeSuccessEvent,
		plugin_dtos.PrepareChatSessionFailedEvent,
		plugin_dtos.PrepareChatSessionSuccessEvent,
		plugin_dtos.HandshakeFailedEvent,
		plugin_dtos.SendChatMessageSuccessEvent,
		plugin_dtos.SendChatMessageFailedEvent,
		plugin_dtos.ResetChatChannelSuccessEvent,
		plugin_dtos.ResetChatChannelFailedEvent:

		p.mu.Lock()
		if ch, exists := p.pendingRequests[payload.RequestID]; exists {
			ch <- payload
		}
		p.mu.Unlock()
		return
	case websocket.ReceivedChatMessageEvent:
		d, _ := GetData[websocket.ReceivedChatMessageEventPayload](payload)
		p.OnReceivedChatMessage(d.SessionID, d.ChannelID, d.Message)
	default:
		p.Logger.Warn("Unknown event", "event", payload.Event)
	}
}

func NewPlugin() *PluginClient {

	logPath := os.Getenv("NAGARE_PLUGIN_LOG_FILE")

	fileRotator := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, fileRotator)
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	p := &PluginClient{
		Config:                &Config{},
		pendingRequests:       make(map[string]chan *websocket.Payload[any]),
		OnReceivedChatMessage: func(sessionID string, _ string, chunk string) {},
		Logger:                slog.New(handler),
	}

	p.connector = NewConnector(p.handleEvent)
	return p
}

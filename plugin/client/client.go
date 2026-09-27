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
	mu                    sync.Mutex
	Metadata              *metadata.PluginMetadata
	ConnectConfig         *ConnectConfig
	ConfigDir             string
	Logger                *slog.Logger
	connector             *Connector
	pendingRequests       map[string]chan *websocket.Payload[any]
	tools                 []PluginTool
	OnReceivedChatMessage func(sessionID string, channelID string, chunk string)
}

func (p *PluginClient) Start(ctx context.Context, onStart func()) error {
	p.ConnectConfig.Port = os.Getenv("NAGARE_PLUGIN_HOST_PORT")
	p.ConnectConfig.ConnectCode = os.Getenv("NAGARE_PLUGIN_CONNECT_CODE")

	err := p.connector.Connect(ctx, fmt.Sprintf("ws://127.0.0.1:%s/ws", p.ConnectConfig.Port))
	if err != nil {
		p.Logger.Error("Start plugin error", "error", err)
		return err
	}

	onStart()
	return nil
}

func (p *PluginClient) handleEvent(payload *websocket.Payload[any]) {
	switch payload.Event {
	case
		// Handshake
		plugin_dtos.HandshakeSuccessEvent,
		plugin_dtos.HandshakeFailedEvent,

		// Prepare chat
		plugin_dtos.PrepareChatSessionFailedEvent,
		plugin_dtos.PrepareChatSessionSuccessEvent,

		// Register plugin tool
		plugin_dtos.RegisterPluginToolSuccessEvent,
		plugin_dtos.RegisterPluginToolFailedEvent,

		// Send chat message
		plugin_dtos.SendChatMessageSuccessEvent,
		plugin_dtos.SendChatMessageFailedEvent,

		// Reset chat channel
		plugin_dtos.ResetChatChannelSuccessEvent,
		plugin_dtos.ResetChatChannelFailedEvent,

		// Register tool categories
		plugin_dtos.RegisterToolCategoriesSuccessEvent,
		plugin_dtos.RegisterToolCategoriesFailedEvent:

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
		ConnectConfig:         &ConnectConfig{},
		pendingRequests:       make(map[string]chan *websocket.Payload[any]),
		OnReceivedChatMessage: func(sessionID string, _ string, chunk string) {},
		Logger:                slog.New(handler),
		tools:                 make([]PluginTool, 0),
	}

	p.ConfigDir = os.Getenv("NAGARE_PLUGIN_CONFIG_DIR")
	p.connector = NewConnector(p.handleEvent)
	return p
}

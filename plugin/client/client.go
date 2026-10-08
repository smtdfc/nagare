package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
	"github.com/smtdfc/nagare/dtos/websocket"
	"github.com/smtdfc/nagare/pkgs/ipc"
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
	if p.ConnectConfig.SocketPath == "" {
		p.ConnectConfig.SocketPath = os.Getenv("NAGARE_PLUGIN_SOCKET_PATH")
	}
	if p.ConnectConfig.SocketPath == "" {
		p.ConnectConfig.SocketPath = ipc.GetDefaultSocketPath()
	}
	if p.ConnectConfig.ConnectCode == "" {
		p.ConnectConfig.ConnectCode = os.Getenv("NAGARE_PLUGIN_CONNECT_CODE")
	}

	err := p.connector.Connect(ctx, p.ConnectConfig.SocketPath)
	if err != nil {
		p.Logger.Error("Start plugin error", "error", err)
		return err
	}

	onStart()
	return p.connector.Wait(ctx)
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
	case plugin_dtos.PluginToolCallEvent:
		go p.handlePluginToolCall(payload)
		return
	case websocket.ReceivedChatMessageEvent:
		d, err := GetData[websocket.ReceivedChatMessageEventPayload](payload)
		if err != nil || d == nil {
			p.Logger.Warn("failed to parse ReceivedChatMessageEvent", "error", err)
			return
		}
		p.OnReceivedChatMessage(d.SessionID, d.ChannelID, d.Message)

	default:
		p.Logger.Warn("Unknown event", "event", payload.Event)
	}
}

func (p *PluginClient) handlePluginToolCall(payload *websocket.Payload[any]) {
	call, err := GetData[plugin_dtos.PluginToolCallEventPayload](payload)
	if err != nil {
		p.sendPluginToolResult(payload.RequestID, "", err.Error())
		return
	}

	p.mu.Lock()
	var calledTool PluginTool
	for _, candidate := range p.tools {
		if candidate.GetName() == call.Name {
			calledTool = candidate
			break
		}
	}
	p.mu.Unlock()

	if calledTool == nil {
		p.sendPluginToolResult(payload.RequestID, "", fmt.Sprintf("plugin tool %q is not registered", call.Name))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	result, err := calledTool.Execute(&ctx, call.Args)
	if err != nil {
		p.sendPluginToolResult(payload.RequestID, "", err.Error())
		return
	}
	p.sendPluginToolResult(payload.RequestID, result, "")
}

func (p *PluginClient) sendPluginToolResult(requestID, result, callError string) {
	err := p.connector.Send(plugin_dtos.PluginToolCallResultEvent, plugin_dtos.PluginToolCallResultEventPayload{
		Result: result,
		Error:  callError,
	}, requestID)
	if err != nil {
		p.Logger.Error("failed to send plugin tool result", "error", err, "request_id", requestID)
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

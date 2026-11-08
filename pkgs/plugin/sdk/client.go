package sdk

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/smtdfc/nagare/pkgs/ipc"
	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
	"github.com/smtdfc/nagare/pkgs/plugin/metadata"
	"gopkg.in/natefinch/lumberjack.v2"
)

type PluginClient struct {
	mu                    sync.Mutex
	Metadata              *metadata.PluginMetadata
	ConnectConfig         *ConnectConfig
	ConfigDir             string
	DataDir               string
	Logger                *slog.Logger
	connector             *Connector
	pendingRequests       map[string]chan *contracts.Payload[any]
	tools                 []PluginTool
	OnReceivedChatMessage func(channelID string, chunk string)
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

func (p *PluginClient) handleEvent(payload *contracts.Payload[any]) {
	switch payload.Event {
	case
		// Handshake
		contracts.HandshakeSuccessEvent,
		contracts.HandshakeFailedEvent,

		// Prepare chat channel
		contracts.PrepareChatChannelFailedEvent,
		contracts.PrepareChatChannelSuccessEvent,

		// Register plugin tool
		contracts.RegisterPluginToolSuccessEvent,
		contracts.RegisterPluginToolFailedEvent,

		// Send chat message
		contracts.SendChatMessageSuccessEvent,
		contracts.SendChatMessageFailedEvent,

		// Reset chat channel
		contracts.ResetChatChannelSuccessEvent,
		contracts.ResetChatChannelFailedEvent,

		// Human-in-the-loop tool approval
		contracts.SubmitToolApprovalSuccessEvent,
		contracts.SubmitToolApprovalFailedEvent,

		// Get model providers
		contracts.GetModelProvidersSuccessEvent,
		contracts.GetModelProvidersFailedEvent,

		// Register chat listener
		contracts.RegisterChatListenerSuccessEvent,
		contracts.RegisterChatListenerFailedEvent,

		// Register tool categories
		contracts.RegisterToolCategoriesSuccessEvent,
		contracts.RegisterToolCategoriesFailedEvent,

		// System prompt injection
		contracts.InjectSystemPromptSuccessEvent,
		contracts.InjectSystemPromptFailedEvent,
		contracts.ClearSystemPromptSuccessEvent,
		contracts.ClearSystemPromptFailedEvent:

		p.mu.Lock()
		if ch, exists := p.pendingRequests[payload.RequestID]; exists {
			ch <- payload
		}
		p.mu.Unlock()
		return
	case contracts.PluginToolCallEvent:
		go p.handlePluginToolCall(payload)
		return
	case contracts.ReceivedChatMessageEvent:
		d, err := GetData[contracts.ReceivedChatMessageEventPayload](payload)
		if err != nil || d == nil {
			p.Logger.Warn("failed to parse ReceivedChatMessageEvent", "error", err)
			return
		}
		p.OnReceivedChatMessage(d.ChannelID, d.Message)

	default:
		p.Logger.Warn("Unknown event", "event", payload.Event)
	}
}

func (p *PluginClient) handlePluginToolCall(payload *contracts.Payload[any]) {
	call, err := GetData[contracts.PluginToolCallEventPayload](payload)
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
	err := p.connector.Send(contracts.PluginToolCallResultEvent, contracts.PluginToolCallResultEventPayload{
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
		pendingRequests:       make(map[string]chan *contracts.Payload[any]),
		OnReceivedChatMessage: func(channelID string, chunk string) {},
		Logger:                slog.New(handler),
		tools:                 make([]PluginTool, 0),
	}

	p.ConfigDir = os.Getenv("NAGARE_PLUGIN_CONFIG_DIR")
	p.DataDir = os.Getenv("NAGARE_PLUGIN_DATA_DIR")
	p.connector = NewConnector(p.handleEvent)
	return p
}

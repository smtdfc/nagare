package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/paths"
	"github.com/smtdfc/nagare/shared/ws"
)

type EventHandler func(msg *dto.WsMessage[any])

func GetServerConnectCodeAddress() (string, error) {
	addr, isExist := os.LookupEnv("NAGARE_PLUGIN_CONNECT_CODE")
	if !isExist {
		return "", fmt.Errorf("NAGARE_PLUGIN_CONNECT_CODE not set")
	}

	return addr, nil
}

type PluginClient struct {
	Logger     *slog.Logger
	serverAddr string
	name       string
	httpClient *req.Client
	ws         *ws.WsInstance

	mu       sync.RWMutex
	handlers map[dto.WsEvent][]EventHandler
}

func (p *PluginClient) CheckServerHealth() error {
	p.Logger.Info("Checking server health", "serverAddr", p.serverAddr)
	_, err := p.httpClient.R().Get(fmt.Sprintf("%s/api/v1/health/check", p.serverAddr))
	if err != nil {
		p.Logger.Error("Server health check failed", "serverAddr", p.serverAddr, "error", err)
		return err
	}

	p.Logger.Info("Server health check passed", "serverAddr", p.serverAddr)
	return nil
}

func (p *PluginClient) On(event dto.WsEvent, handler EventHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[event] = append(p.handlers[event], handler)
	p.Logger.Debug("Registered event listener", "event", event)
}

func (p *PluginClient) Off(event dto.WsEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.handlers, event)
	p.Logger.Debug("Removed all event listeners for", "event", event)
}

func (p *PluginClient) dispatchEvent(msg *dto.WsMessage[any]) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if handlers, exists := p.handlers[msg.Event]; exists {
		for _, handler := range handlers {
			go handler(msg)
		}
	}
}

func (p *PluginClient) Start() error {
	p.Logger.Info("Starting plugin", "name", p.name)
	serverAddr, err := helpers.GetRestApiConnect()
	if err != nil {
		p.Logger.Error("Failed to get address", "error", err)
		return err
	}

	wsAddr, err := helpers.GetPluginWebsocketConnect()
	if err != nil {
		p.Logger.Error("Failed to get websocket address", "error", err)
		return err
	}

	connectCode, err := GetServerConnectCodeAddress()
	if err != nil {
		p.Logger.Error("Failed to get connect code", "error", err)
		return err
	}

	p.serverAddr = serverAddr
	if err := p.CheckServerHealth(); err != nil {
		return err
	}

	c, err := ws.Dial(wsAddr)
	if err != nil {
		p.Logger.Error("Failed to connect to server", "error", err)
		panic("Failed to connect to server: " + err.Error())
	}

	p.ws = c
	defer p.ws.Close()

	err = ws.SendMessage(c, dto.WS_PLUGIN_REGISTER, dto.PluginRegisterEvent{
		Code: connectCode,
	})
	if err != nil {
		p.Logger.Error("Failed to register plugin", "error", err)
		return err
	}

	for {
		msg, err := p.ws.ReadMessage()
		if err != nil {
			p.Logger.Error("Failed to read message / connection dropped", "error", err)
			break
		}

		p.Logger.Info("Received message from server", "event", msg.Event)

		p.dispatchEvent(msg)
	}
	return nil
}

func NewPluginClient(name string) *PluginClient {
	logDir := path.Join(paths.LogDir, "plugins", name)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	timestamp := time.Now().Format("2006-01-02T15:04:05")
	file, err := os.OpenFile(path.Join(logDir, timestamp+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	jsonHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(jsonHandler)

	return &PluginClient{
		Logger:     logger,
		name:       name,
		httpClient: req.C().SetCommonHeader("Accept", "application/json").SetCommonHeader("X-nagare-plugin", name),
		handlers:   make(map[dto.WsEvent][]EventHandler),
	}
}

func WsRequest[TResponse any, TError any](
	p *PluginClient,
	sendEvent dto.WsEvent,
	successEvent dto.WsEvent,
	failureEvent dto.WsEvent,
	payload any,
	timeoutMs time.Duration,
) (*TResponse, *TError, error) {
	if timeoutMs == 0 {
		timeoutMs = 10 * time.Second
	}

	respChan := make(chan *TResponse, 1)
	errChan := make(chan *TError, 1)

	var successHandler EventHandler
	successHandler = func(msg *dto.WsMessage[any]) {
		p.Off(successEvent)
		p.Off(failureEvent)

		var target TResponse
		bytesData, err := json.Marshal(msg.Payload)
		if err != nil {
			p.Logger.Error("Failed to marshal ws response data", "error", err)
			return
		}
		if err := json.Unmarshal(bytesData, &target); err != nil {
			p.Logger.Error("Failed to unmarshal ws response data", "error", err)
			return
		}

		respChan <- &target
	}

	var failureHandler EventHandler
	failureHandler = func(msg *dto.WsMessage[any]) {
		p.Off(successEvent)
		p.Off(failureEvent)

		var target TError
		bytesData, err := json.Marshal(msg.Payload)
		if err != nil {
			p.Logger.Error("Failed to marshal ws error data", "error", err)
			return
		}
		if err := json.Unmarshal(bytesData, &target); err != nil {
			p.Logger.Error("Failed to unmarshal ws error data", "error", err)
			return
		}

		errChan <- &target
	}

	p.On(successEvent, successHandler)
	p.On(failureEvent, failureHandler)

	err := ws.SendMessage(p.ws, sendEvent, payload)
	if err != nil {
		p.Off(successEvent)
		p.Off(failureEvent)
		return nil, nil, fmt.Errorf("failed to send ws request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutMs)
	defer cancel()

	select {
	case <-ctx.Done():
		p.Off(successEvent)
		p.Off(failureEvent)
		return nil, nil, fmt.Errorf("websocket request timeout for event: %v", sendEvent)

	case res := <-respChan:
		return res, nil, nil

	case terr := <-errChan:
		return nil, terr, nil
	}
}

package client

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/imroc/req/v3"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/paths"
	"github.com/smtdfc/nagare/shared/ws"
)

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

func (p *PluginClient) Handshake() error {
	return helpers.Unimplemented()
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
			p.Logger.Error("Failed to read message", "error", err)
			break
		}
		p.Logger.Error("Received message from server", "event", msg.Event)
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
	}
}

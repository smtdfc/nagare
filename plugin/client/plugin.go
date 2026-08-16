package client

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/imroc/req/v3"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/paths"
)

func GetPort() (string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return "", errors.New("no port specified")
	}
	return args[0], nil
}

type PluginClient struct {
	Logger     *slog.Logger
	serverAddr string
	name       string
	httpClient *req.Client
}

func (p *PluginClient) CheckServerHealth() error {
	p.Logger.Info("Checking server health", "serverAddr", p.serverAddr)
	_, err := p.httpClient.R().Get(fmt.Sprintf("%s/api/v1/health/check", p.serverAddr))
	if err != nil {
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
	port, err := GetPort()
	if err != nil {
		return err
	}

	serverAddr := "http://localhost:" + port
	p.serverAddr = serverAddr
	if err := p.CheckServerHealth(); err != nil {
		return err
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

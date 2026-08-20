package client

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/paths"
	"github.com/smtdfc/nagare/shared/plugin"
	"github.com/smtdfc/nagare/shared/ws"
)

type EventHandler func(msg *dto.WsMessage[any])

type PluginClient struct {
	Metadata   *plugin.PluginMetadata
	Logger     *slog.Logger
	serverAddr string
	name       string
	httpClient *req.Client
	ws         *ws.WsInstance

	mu       sync.RWMutex
	handlers map[dto.WsEvent]map[string]EventHandler
}

func (p *PluginClient) Start(ctx context.Context, onReady func(context.Context) error) error {
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

	type readResult struct {
		msg *dto.WsMessage[any]
		err error
	}
	msgChan := make(chan readResult)

	go func() {
		for {
			msg, err := p.ws.ReadMessage()
			msgChan <- readResult{msg: msg, err: err}
			if err != nil {
				break
			}
		}
	}()

	if onReady != nil {
		go func() {
			if err := onReady(ctx); err != nil {
				p.Logger.Error("onReady callback failed", "error", err)
			}
		}()
	}

	for {
		select {
		case <-ctx.Done():
			p.Logger.Info("Context cancelled, stopping plugin client...", "reason", ctx.Err())
			return ctx.Err()

		case res := <-msgChan:
			if res.err != nil {
				if errors.Is(res.err, io.EOF) {
					p.Logger.Info("Stream session ended gracefully by server.")
					continue
				}

				p.Logger.Error("Failed to read message / connection dropped", "error", res.err)
				p.mu.Lock()
				if p.ws != nil {
					p.ws.Close()
					p.ws = nil
				}
				p.mu.Unlock()
				return nil
			}
			p.dispatchEvent(res.msg)
		}
	}
}

func (p *PluginClient) Register() error {
	connectCode, err := GetServerConnectCodeAddress()
	if err != nil {
		p.Logger.Error("Failed to get connect code", "error", err)
		return err
	}

	_, werr, err := WsRequest[dto.PluginRegisterSuccessEvent, dto.PluginRegisterFailedEvent](
		p,
		dto.WS_PLUGIN_REGISTER,
		dto.WS_PLUGIN_REGISTER_SUCCESS,
		dto.WS_PLUGIN_REGISTER_FAILED,
		&dto.PluginRegisterEvent{
			Code: connectCode,
		},
		2*time.Minute,
	)

	if werr != nil {
		return errors.New(werr.Cause)
	}

	if err != nil {
		return err

	}

	return nil
}

func NewPluginClient(name string) *PluginClient {
	logDir := filepath.Join(paths.LogDir, "plugins", name)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	timestamp := time.Now().Format("2006-01-02T15:04:05")
	file, err := os.OpenFile(filepath.Join(logDir, timestamp+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	jsonHandler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(jsonHandler)

	return &PluginClient{
		Logger:     logger.With("plugin", name),
		name:       name,
		httpClient: req.C().SetCommonHeader("Accept", "application/json").SetCommonHeader("X-nagare-plugin", name),
		handlers:   make(map[dto.WsEvent]map[string]EventHandler),
	}
}

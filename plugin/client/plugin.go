package client

import (
	"errors"
	"fmt"
	"os"

	"github.com/imroc/req/v3"
	"github.com/smtdfc/nagare/shared/helpers"
)

func GetPort() (string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return "", errors.New("no port specified")
	}
	return args[0], nil
}

type PluginClient struct {
	serverAddr string
	name       string
	httpClient *req.Client
}

func (p *PluginClient) CheckServerHealth() error {
	_, err := p.httpClient.R().Get(fmt.Sprintf("%s/api/v1/health/check", p.serverAddr))
	if err != nil {
		return err
	}
	return nil
}

func (p *PluginClient) Handshake() error {
	return helpers.Unimplemented()
}

func (p *PluginClient) Start() error {
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
	return &PluginClient{
		name:       name,
		httpClient: req.C().SetCommonHeader("Accept", "application/json").SetCommonHeader("X-nagare-plugin", name),
	}
}

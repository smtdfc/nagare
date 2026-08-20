package client

import (
	"fmt"
	"os"
)

func GetServerConnectCodeAddress() (string, error) {
	addr, isExist := os.LookupEnv("NAGARE_PLUGIN_CONNECT_CODE")
	if !isExist {
		return "", fmt.Errorf("NAGARE_PLUGIN_CONNECT_CODE not set")
	}

	return addr, nil
}

func (p *PluginClient) CheckServerHealth() error {
	p.Logger.Info("Checking server health", "serverAddr", p.serverAddr)
	_, err := p.httpClient.R().Get(fmt.Sprintf("%s/health/check", p.serverAddr))
	if err != nil {
		p.Logger.Error("Server health check failed", "serverAddr", p.serverAddr, "error", err)
		return err
	}

	p.Logger.Info("Server health check passed", "serverAddr", p.serverAddr)
	return nil
}

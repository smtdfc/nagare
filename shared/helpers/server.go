package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func CheckServerRun() (bool, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	serverHost, err := GetRestApiConnect()
	if err != nil {
		return false, err
	}
	url := fmt.Sprintf("%s/health/check", serverHost)

	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, fmt.Errorf("server error: %d", resp.StatusCode)
}

func GetWebsocketConnect() (string, error) {
	serverHost, err := ResolveServerHost()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("ws://%s/ws/chat", serverHost)
	return url, nil
}

func GetPluginWebsocketConnect() (string, error) {
	serverHost, err := ResolveServerHost()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("ws://%s/ws/plugin", serverHost)
	return url, nil
}

func GetRestApiConnect() (string, error) {
	serverHost, err := ResolveServerHost()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("http://%s/api/v1", serverHost)
	return url, nil
}

func GetServerPort() (string, error) {
	osPort := os.Getenv("NAGARE_SERVER_PORT")
	if osPort != "" {
		return osPort, nil
	}

	return "9832", nil
}

func IsRunWithRemoteServer() bool {
	osPort := os.Getenv("NAGARE_SERVER_HOST")
	return osPort != ""
}

func GetServerRemoteUrl() string {
	osHost := os.Getenv("NAGARE_SERVER_HOST")
	return osHost
}

func ResolveServerHost() (string, error) {
	if IsRunWithRemoteServer() {
		return GetServerRemoteUrl(), nil
	}
	port, err := GetServerPort()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("localhost:%s", port), nil
}

func GetServerPublicKey() string {
	return os.Getenv("NAGARE_PUBLIC_KEY")
}

package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"

	"github.com/smtdfc/nagare/shared/paths"
)

func CheckServerRun(port string) (bool, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	url := fmt.Sprintf("http://localhost:%s/api/v1/health/check", port)

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

func TryStartServer(port string) error {
	cmd := exec.Command(paths.ServerBinFile, "--port", port)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Cannot start server. Cause: %w", err)
	}

	return nil
}

func GetWebsocketConnect(port string) (string, error) {
	url := fmt.Sprintf("http://localhost:%s/ws/chat", port)
	return url, nil
}

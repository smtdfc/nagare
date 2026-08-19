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

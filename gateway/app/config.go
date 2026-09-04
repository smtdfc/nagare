package app

import (
	"fmt"
	"os"
)

type AppConfig struct {
	Port      string
	DebugMode bool
	PublicKey string
}

// @Injectable
func ResolveConfig() (*AppConfig, error) {
	conf := &AppConfig{
		Port: "9832",
	}

	if port := os.Getenv("NAGARE_GATEWAY_PORT"); port != "" {
		conf.Port = port
	}

	publicKey := os.Getenv("NAGARE_GATEWAY_PUBLIC_KEY")
	if publicKey == "" {
		return nil, fmt.Errorf("missing critical configuration: NAGARE_GATEWAY_PUBLIC_KEY is required")
	}
	conf.PublicKey = publicKey

	if os.Getenv("NAGARE_GATEWAY_MODE") == "debug" {
		conf.DebugMode = true
	}

	return conf, nil
}

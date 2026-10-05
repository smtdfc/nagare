package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port      string
	Host      string
	DebugMode bool
	PublicKey string
}

// @Injectable
func ResolveConfig() (*Config, error) {
	conf := &Config{
		Port: "9832",
	}

	if host := os.Getenv("NAGARE_GATEWAY_HOST"); host != "" {
		conf.Host = host
	} else {
		conf.Host = "localhost"
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

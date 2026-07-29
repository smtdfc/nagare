package app

import (
	"flag"

	"github.com/smtdfc/nagare/server/config"
)

type CliInjectedConfiguration struct {
	Port      string
	PublicKey string
}

// @Injectable
func NewCliInjectedConfiguration() *CliInjectedConfiguration {
	port := flag.String("port", "3000", "Port")
	pubKey := flag.String("pubkey", "", "RSA Public Key base64")

	flag.Parse()
	return &CliInjectedConfiguration{
		Port:      *port,
		PublicKey: *pubKey,
	}
}

// @Injectable
func NewConfig(cliInjected *CliInjectedConfiguration) *config.ServerConfig {
	defaultConfig := &config.ServerConfig{
		Port:      cliInjected.Port,
		PublicKey: cliInjected.PublicKey,
	}
	return defaultConfig
}

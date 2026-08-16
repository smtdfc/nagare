package app

import (
	"flag"

	"github.com/smtdfc/nagare/server/config"
	"github.com/smtdfc/nagare/shared/helpers"
)

type CliInjectedConfiguration struct {
	Port string
}

// @Injectable
func NewCliInjectedConfiguration() *CliInjectedConfiguration {
	port := flag.String("port", "3000", "Port")

	flag.Parse()
	return &CliInjectedConfiguration{
		Port: *port,
	}
}

// @Injectable
func NewConfig(cliInjected *CliInjectedConfiguration) *config.ServerConfig {
	defaultConfig := &config.ServerConfig{
		Port:      cliInjected.Port,
		PublicKey: helpers.GetServerPublicKey(),
	}
	return defaultConfig
}

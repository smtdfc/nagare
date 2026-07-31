package config

import (
	"github.com/smtdfc/nagare/core/domains"
	persistence "github.com/smtdfc/nagare/core/persistence/config"
)

type ConfigManager struct {
	io *persistence.ConfigIO
}

func (c *ConfigManager) Load() (*domains.Config, error) {
	return c.io.Read()
}

func (c *ConfigManager) Save(conf *domains.Config) error {
	return c.io.Write(conf)
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		io: persistence.NewConfigIO(),
	}
}

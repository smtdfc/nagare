package config

type ServerConfig struct {
	PublicKey         string `json:"public_key"`
	Port              string `json:"port"`
	ReduceMemoryUsage bool
}

func (c *ServerConfig) Merge(config *ServerConfig) {
	if config == nil {
		return
	}

	if config.PublicKey != "" {
		c.PublicKey = config.PublicKey
	}

	if config.Port != "" {
		c.Port = config.Port
	}
}

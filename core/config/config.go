package config

import "github.com/smtdfc/nagare/core/domains"

func GetListProvider(c *domains.Config) []domains.ProviderConfig {
	list := []domains.ProviderConfig{}
	for _, provider := range c.Providers {
		list = append(list, *provider)
	}
	return list
}

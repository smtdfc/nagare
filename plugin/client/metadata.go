package client

import (
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/plugin"
)

func (p *PluginClient) LoadMetadataFromString(raw string) error {
	metadata, err := helpers.FromJson[plugin.PluginMetadata](raw)
	if err != nil {
		return err
	}

	p.Metadata = metadata
	return nil
}

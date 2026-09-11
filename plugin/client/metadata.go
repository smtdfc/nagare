package client

import (
	mt "github.com/smtdfc/nagare/plugin/metadata"
	"github.com/smtdfc/nagare/shared/helpers"
)

func (p *PluginClient) LoadMetadata(raw string) (*mt.PluginMetadata, error) {
	metadata, err := helpers.UnmarshalJson[mt.PluginMetadata](raw)
	if err != nil {
		return nil, err
	}

	p.Metadata = metadata
	return metadata, nil
}

package manager

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/discovery"
	"github.com/smtdfc/nagare/core/discovery/providers"
	"github.com/smtdfc/nagare/core/logger"
)

type DiscoveryManager struct {
	providers map[string]discovery.SourceProviderAdapter
	logger    *logger.BaseLogger
}

func (d *DiscoveryManager) GetAvailableSourceProvider() []string {
	result := make([]string, 0, len(d.providers))
	for name, _ := range d.providers {
		result = append(result, name)
	}

	return result
}

// AddProvider registers or replaces a discovery source provider.
func (d *DiscoveryManager) AddProvider(name string, provider discovery.SourceProviderAdapter) {
	if d.providers == nil {
		d.providers = make(map[string]discovery.SourceProviderAdapter)
	}

	d.providers[name] = provider
}

func (d *DiscoveryManager) Search(ctx context.Context, provider string, query string, limit int) ([]discovery.Source, error) {
	adapter, exist := d.providers[provider]
	if !exist {
		return nil, custom_errors.ErrSourceProviderNotFound
	}

	result, err := adapter.Search(ctx, query, limit)
	if err != nil {
		return nil, custom_errors.ErrDiscoveryFailed
	}

	return result, nil
}

func (d *DiscoveryManager) Fetch(ctx context.Context, provider string, id string) (*discovery.Detail, error) {
	adapter, exist := d.providers[provider]
	if !exist {
		return nil, custom_errors.ErrSourceProviderNotFound
	}

	result, err := adapter.Fetch(ctx, id)
	if err != nil {
		return nil, custom_errors.ErrDiscoveryFailed
	}

	return result, nil
}

// @Injectable
func NewDiscoveryManager(logger *logger.BaseLogger) *DiscoveryManager {
	return &DiscoveryManager{
		providers: map[string]discovery.SourceProviderAdapter{
			"Wikipedia": providers.NewWikipediaProvider(),
		},
		logger: logger.With("module", "discovery-manager"),
	}
}

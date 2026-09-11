package plugin

import (
	"context"

	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/shared/dtos/rest"
	"github.com/smtdfc/nagare/shared/helpers"
)

func toPluginDTO(domain *plugin.Plugin) *rest.Plugin {
	features := make([]string, 0, len(domain.Features))
	for _, f := range domain.Features {
		features = append(features, f.ToString())
	}

	return &rest.Plugin{
		ID:       domain.ID.String(),
		PluginID: domain.PluginID,
		Name:     domain.Name,
		Author:   domain.Author,
		Features: features,
		Version:  domain.Version,
		IsActive: domain.IsActive,
	}
}

type Service struct {
	pluginMgr *manager.PluginManager
}

func (p *Service) ListPlugins(ctx context.Context) (*rest.GetListPluginResponse, error) {
	plugins, err := p.pluginMgr.GetListPlugin(ctx)
	if err != nil {
		return nil, err
	}

	return &rest.GetListPluginResponse{
		Plugins: helpers.SliceMap(plugins, toPluginDTO),
	}, nil
}

func (p *Service) InstallLocalPlugin(ctx context.Context, request *rest.InstallLocalPluginRequest) error {
	return p.pluginMgr.Install(ctx, request.Path)
}

// @Injectable
func NewPluginService(
	pluginMgr *manager.PluginManager,
) *Service {
	return &Service{
		pluginMgr: pluginMgr,
	}
}

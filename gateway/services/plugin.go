package services

import (
	"context"

	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/shared/dtos/rest"
	"github.com/smtdfc/nagare/shared/helpers"
)

func toPluginDTO(domain *plugin.Plugin) *rest.Plugin {
	return &rest.Plugin{
		ID:       "",
		PluginID: "",
		Name:     "",
		Author:   "",
		Features: []string{},
		Version:  "",
		IsActive: false,
	}
}

type PluginService struct {
	pluginMgr *manager.PluginManager
}

func (p *PluginService) GetListPlugin(ctx context.Context) (*rest.GetListPluginResponse, error) {
	plugins, err := p.pluginMgr.GetListPlugin(ctx)
	if err != nil {
		return nil, err
	}

	return &rest.GetListPluginResponse{
		Plugins: helpers.SliceMap(plugins, toPluginDTO),
	}, nil
}

func (p *PluginService) InstallLocalPlugin(ctx context.Context, request *rest.InstallLocalPluginRequest) error {
	return p.pluginMgr.Install(ctx, request.Path)
}

// @Injectable
func NewPluginService(
	pluginMgr *manager.PluginManager,
) *PluginService {
	return &PluginService{
		pluginMgr: pluginMgr,
	}
}

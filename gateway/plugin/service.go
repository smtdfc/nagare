package plugin

import (
	"context"

	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/dtos/rest"
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

type PluginService struct {
	pluginMgr *manager.PluginManager
}

func (p *PluginService) ListPlugins(ctx context.Context) (*rest.GetListPluginResponse, error) {
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

func (p *PluginService) UninstallPlugin(ctx context.Context, request *rest.UninstallPluginRequest) error {
	return p.pluginMgr.Uninstall(ctx, request.ID)
}

func (p *PluginService) ActivatePlugin(ctx context.Context, request *rest.ActivatePluginRequest) error {
	return p.pluginMgr.Activate(ctx, request.ID)
}

func (p *PluginService) DeactivatePlugin(ctx context.Context, request *rest.DeactivatePluginRequest) error {
	return p.pluginMgr.Deactivate(ctx, request.ID)
}

func (p *PluginService) GetPluginStatus(ctx context.Context, request *rest.GetPluginStatusRequest) (*rest.GetPluginStatusResponse, error) {
	status, err := p.pluginMgr.GetPluginStatus(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &rest.GetPluginStatusResponse{
		Status: &rest.PluginStatus{
			PID:         status.PID,
			PluginID:    status.PluginID,
			Name:        status.Name,
			Version:     status.Version,
			CPUPercent:  status.CPUPercent,
			MemoryUsage: status.MemoryUsage,
		},
	}, nil
}

// @Injectable
func NewPluginService(
	pluginMgr *manager.PluginManager,
) *PluginService {
	return &PluginService{
		pluginMgr: pluginMgr,
	}
}

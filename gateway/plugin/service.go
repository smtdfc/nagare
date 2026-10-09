package plugin

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/media/upload"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/pkgs/helpers"
)

func toPluginDTO(domain *plugin.Plugin) *rest.Plugin {
	features := make([]string, 0, len(domain.Features))
	for _, f := range domain.Features {
		features = append(features, f.ToString())
	}

	return &rest.Plugin{
		ID:          domain.ID.String(),
		PackageName: domain.PackageName,
		Name:        domain.Name,
		Author:      domain.Author,
		Features:    features,
		Version:     domain.Version,
		IsActive:    domain.IsActive,
	}
}

type Service struct {
	pluginMgr *manager.PluginManager
	uploadMgr *upload.Manager
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

func (p *Service) UploadPlugin(_ context.Context, savePath string) (*rest.UploadPluginResponse, error) {
	id, err := p.uploadMgr.AddAttachment(savePath)
	if err != nil {
		return nil, err
	}

	return &rest.UploadPluginResponse{AttachmentID: id}, nil
}

func (p *Service) InstallLocalPlugin(ctx context.Context, request *rest.InstallLocalPluginRequest) (*rest.InstallLocalPluginResponse, error) {
	plg, err := p.pluginMgr.Install(ctx, request.Path)
	if err != nil {
		return nil, err
	}

	return &rest.InstallLocalPluginResponse{Plugin: toPluginDTO(plg)}, nil
}

func (p *Service) UninstallPlugin(ctx context.Context, request *rest.UninstallPluginRequest) error {
	return p.pluginMgr.Uninstall(ctx, request.ID)
}

func (p *Service) ActivatePlugin(ctx context.Context, request *rest.ActivatePluginRequest) error {
	return p.pluginMgr.Activate(ctx, request.ID)
}

func (p *Service) DeactivatePlugin(ctx context.Context, request *rest.DeactivatePluginRequest) error {
	return p.pluginMgr.Deactivate(ctx, request.ID)
}

func (p *Service) GetPluginStatus(ctx context.Context, request *rest.GetPluginStatusRequest) (*rest.GetPluginStatusResponse, error) {
	status, err := p.pluginMgr.GetPluginStatus(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &rest.GetPluginStatusResponse{
		Status: &rest.PluginStatus{
			PID:         status.PID,
			PackageName: status.PackageName,
			Name:        status.Name,
			Version:     status.Version,
			CPUPercent:  status.CPUPercent,
			MemoryUsage: status.MemoryUsage,
		},
	}, nil
}

func (p *Service) InstallPluginFromAttachment(ctx context.Context, request *rest.InstallPluginFromAttachmentRequest) (*rest.InstallPluginFromAttachmentResponse, error) {
	attachmentPath, exists := p.uploadMgr.GetAttachmentPath(request.AttachmentID)
	if !exists {
		return nil, custom_errors.ErrPluginNotFound
	}

	plg, err := p.pluginMgr.Install(ctx, attachmentPath)
	if err != nil {
		return nil, err
	}

	return &rest.InstallPluginFromAttachmentResponse{
		Plugin: toPluginDTO(plg),
	}, nil
}

// @Injectable
func NewPluginService(
	pluginMgr *manager.PluginManager,
	uploadMgr *upload.Manager,
) *Service {
	return &Service{
		pluginMgr: pluginMgr,
		uploadMgr: uploadMgr,
	}
}

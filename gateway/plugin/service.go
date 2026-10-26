package plugin

import (
	"context"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/media/upload"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/dtos/rest"
)

func signatureWarning(status string) string {
	switch status {
	case plugin.SignatureStatusUnsigned:
		return "This plugin is not signed. Only install plugins from a source you trust."
	case plugin.SignatureStatusUnverified:
		return "This plugin has a signature, but no trusted public key was found. Add the publisher's public key to the trusted folder and reinstall."
	default:
		return ""
	}
}

type Service struct {
	pluginMgr    *manager.PluginManager
	uploadMgr    *upload.Manager
	pluginMapper IPluginDtoMapper
}

func (p *Service) ListPlugins(ctx context.Context) (*rest.GetListPluginResponse, error) {
	plugins, err := p.pluginMgr.GetListPlugin(ctx)
	if err != nil {
		return nil, err
	}

	return &rest.GetListPluginResponse{
		Plugins: p.pluginMapper.ToPluginDTOs(plugins),
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

	return &rest.InstallLocalPluginResponse{
		Plugin:  p.pluginMapper.ToPluginDTO(plg),
		Warning: signatureWarning(plg.SignatureStatus),
	}, nil
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

func (p *Service) signatureAttachmentPath(signatureAttachmentID string) (string, error) {
	if signatureAttachmentID == "" {
		return "", nil
	}

	path, exists := p.uploadMgr.GetAttachmentPath(signatureAttachmentID)
	if !exists {
		return "", custom_errors.ErrPluginSignatureInvalid
	}

	return path, nil
}

func (p *Service) InstallPluginFromAttachment(ctx context.Context, request *rest.InstallPluginFromAttachmentRequest) (*rest.InstallPluginFromAttachmentResponse, error) {
	attachmentPath, exists := p.uploadMgr.GetAttachmentPath(request.AttachmentID)
	if !exists {
		return nil, custom_errors.ErrPluginNotFound
	}

	signaturePath, err := p.signatureAttachmentPath(request.SignatureAttachmentID)
	if err != nil {
		return nil, err
	}

	plg, err := p.pluginMgr.InstallWithSignature(ctx, attachmentPath, signaturePath)
	if err != nil {
		return nil, err
	}

	return &rest.InstallPluginFromAttachmentResponse{
		Plugin:  p.pluginMapper.ToPluginDTO(plg),
		Warning: signatureWarning(plg.SignatureStatus),
	}, nil
}

func (p *Service) UpdatePlugin(ctx context.Context, request *rest.UpdatePluginRequest) (*rest.UpdatePluginResponse, error) {
	attachmentPath, exists := p.uploadMgr.GetAttachmentPath(request.AttachmentID)
	if !exists {
		return nil, custom_errors.ErrPluginNotFound
	}

	signaturePath, err := p.signatureAttachmentPath(request.SignatureAttachmentID)
	if err != nil {
		return nil, err
	}

	plg, err := p.pluginMgr.UpdateWithSignature(ctx, attachmentPath, signaturePath)
	if err != nil {
		return nil, err
	}

	return &rest.UpdatePluginResponse{
		Plugin:  p.pluginMapper.ToPluginDTO(plg),
		Warning: signatureWarning(plg.SignatureStatus),
	}, nil
}

// @Injectable
func NewPluginService(
	pluginMgr *manager.PluginManager,
	uploadMgr *upload.Manager,
	pluginMapper IPluginDtoMapper,
) *Service {
	return &Service{
		pluginMgr:    pluginMgr,
		uploadMgr:    uploadMgr,
		pluginMapper: pluginMapper,
	}
}

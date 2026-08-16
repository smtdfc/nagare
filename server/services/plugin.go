package services

import (
	"github.com/smtdfc/nagare/core/global"
	"github.com/smtdfc/nagare/server/custom_errors"
	"github.com/smtdfc/nagare/shared/dto"
)

type PluginService struct {
}

func (s *PluginService) GetAllPlugin() (*dto.GetAllPluginsResponse, error) {
	plugins, err := global.GlobalPluginMgr.GetAllPlugins()
	if err != nil {
		return nil, custom_errors.NewServiceError(err.Error(), 400)
	}

	pluginInfos := make([]dto.PluginInfo, 0, len(plugins))
	for _, plugin := range plugins {
		pluginInfos = append(pluginInfos, dto.PluginInfo{
			ID:         plugin.ID,
			PluginID:   plugin.PluginID,
			Name:       plugin.Name,
			Active:     plugin.Active,
			ApiVersion: plugin.ApiVersion,
			Author:     plugin.Author,
			Version:    plugin.Version,
		})
	}
	return &dto.GetAllPluginsResponse{Plugins: pluginInfos}, nil
}

func (s *PluginService) InstallLocalPlugin(req *dto.InstallLocalPluginRequest) (*dto.InstallLocalPluginResponse, error) {
	err := global.GlobalPluginMgr.RegisterPlugin(req.Path)
	if err != nil {
		return nil, custom_errors.NewServiceError(err.Error(), 400)
	}
	return &dto.InstallLocalPluginResponse{}, nil
}

// @Injectable
func NewPluginService() *PluginService {
	return &PluginService{}
}

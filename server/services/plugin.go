package services

import (
	"github.com/smtdfc/nagare/core/global"
	"github.com/smtdfc/nagare/shared/dto"
)

type PluginService struct {
}

func (s *PluginService) GetAllPlugin() (*dto.GetAllPluginsResponse, error) {
	plugins, err := global.GlobalPluginMgr.GetAllPlugins()
	if err != nil {
		return nil, err
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

// @Injectable
func NewPluginService() *PluginService {
	return &PluginService{}
}

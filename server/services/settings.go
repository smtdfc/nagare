package services

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/global"
	"github.com/smtdfc/nagare/shared/dto"
)

type SettingsService struct {
}

func (s *SettingsService) GetGeneralSettings() (*dto.GetGeneralSettingsResponse, error) {
	settings, err := global.GlobalConfigMgr.GetGeneralConfig()
	if err != nil {
		return nil, err
	}
	return &dto.GetGeneralSettingsResponse{
		Settings: &dto.GeneralSettings{
			CurrentModel:    settings.CurrentModel,
			CurrentProvider: settings.CurrentProvider,
		},
	}, nil
}

func (s *SettingsService) SaveGeneralSettings(request *dto.SaveGeneralSettingsRequest) (*dto.SaveGeneralSettingsResponse, error) {
	settings := domains.GeneralConfig{
		CurrentModel:    request.Settings.CurrentModel,
		CurrentProvider: request.Settings.CurrentProvider,
	}
	err := global.GlobalConfigMgr.SaveGeneralConfig(&settings)
	if err != nil {
		return nil, err
	}
	return &dto.SaveGeneralSettingsResponse{}, nil
}

// @Injectable
func NewSettingsService() *SettingsService {
	return &SettingsService{}
}

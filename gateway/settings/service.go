package settings

import (
	"context"
	"log/slog"

	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/config/manager"
	"github.com/smtdfc/nagare/dtos/rest"
	"github.com/smtdfc/nagare/gateway/common/custom_errors"
)

type SettingsService struct {
	configMgr *manager.ConfigManager
	logger    *slog.Logger
}

func (s *SettingsService) SetGeneralSettings(ctx context.Context, request *rest.SetGeneralSettingsRequest) error {
	if request.GeneralSettings == nil {
		return custom_errors.ErrInvalidBody
	}

	generalConfig := &config.GeneralConfig{
		DefaultLLMProvider: request.GeneralSettings.DefaultLLMProvider,
		DefaultLLMModel:    request.GeneralSettings.DefaultLLMModel,
	}

	err := s.configMgr.SetGeneralConfig(ctx, generalConfig)
	if err != nil {
		return err
	}

	return nil
}

func (s *SettingsService) GetGeneralSettings(ctx context.Context) (*rest.GetGeneralSettingsResponse, error) {

	generalConfig, err := s.configMgr.GetGeneralConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &rest.GetGeneralSettingsResponse{
		GeneralSettings: &rest.GeneralSettings{
			DefaultLLMModel:    generalConfig.DefaultLLMModel,
			DefaultLLMProvider: generalConfig.DefaultLLMProvider,
		},
	}, nil
}

// @Injectable
func NewSettingsService(configMgr *manager.ConfigManager) *SettingsService {
	return &SettingsService{
		configMgr: configMgr,
	}
}

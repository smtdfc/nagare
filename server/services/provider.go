package services

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/providers"
	"github.com/smtdfc/nagare/server/custom_errors"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/helpers"
)

type ProviderService struct{}

func (s *ProviderService) GetListProvider() (*dto.GetListProviderResponse, error) {
	providers, err := helpers.Map(providers.GetAllProviderConfig(), func(v domains.ProviderConfig) (dto.Provider, error) {
		return dto.Provider{
			ID:              v.ID,
			Compatible:      string(v.Compatible),
			Name:            v.Name,
			BaseURL:         v.BaseURL,
			APIKey:          v.APIKey,
			IsEnable:        v.IsEnable,
			AvailableModels: v.AvailableModels,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	resp := &dto.GetListProviderResponse{
		Providers: providers,
	}
	return resp, nil
}

func (s *ProviderService) GetProviderDetails(id string) (*dto.GetProviderByIDResponse, error) {
	v, isExist := providers.FindProviderConfigByID(id)
	if !isExist {
		return nil, custom_errors.NewServiceError(
			"Provider not found",
			400,
		)
	}

	resp := &dto.GetProviderByIDResponse{
		Provider: dto.Provider{
			ID:              v.ID,
			Compatible:      string(v.Compatible),
			Name:            v.Name,
			BaseURL:         v.BaseURL,
			APIKey:          v.APIKey,
			IsEnable:        v.IsEnable,
			AvailableModels: v.AvailableModels,
		},
	}
	return resp, nil
}

func (s *ProviderService) UpdateProvider(id string, data dto.UpdateProviderRequest) error {
	v, isExist := providers.FindProviderConfigByID(id)
	if !isExist {
		return custom_errors.NewServiceError(
			"Provider not found",
			400,
		)
	}

}

// @Injectable
func NewProviderService() *ProviderService {
	return &ProviderService{}
}

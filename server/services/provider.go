package services

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/providers"
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

// @Injectable
func NewProviderService() *ProviderService {
	return &ProviderService{}
}

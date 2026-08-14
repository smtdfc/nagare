package services

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/global"

	"github.com/smtdfc/nagare/server/custom_errors"
	"github.com/smtdfc/nagare/shared/dto"
	"github.com/smtdfc/nagare/shared/helpers"
)

type ProviderService struct{}

func (s *ProviderService) GetListProvider() (*dto.GetListProviderResponse, error) {
	providerInfos, err := global.GlobalConfigMgr.GetListProviders()
	if err != nil {
		return nil, err
	}

	providers, err := helpers.Map(providerInfos, func(v *domains.LLMProviderConfigInfo) (dto.ProviderInfo, error) {
		return dto.ProviderInfo{
			ID:              v.ID,
			Compatible:      string(v.Compatible),
			Name:            v.Name,
			BaseURL:         v.BaseURL,
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
	v, err := global.GlobalConfigMgr.GetLLMProviderConfigByID(id)
	if err != nil {
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

func (s *ProviderService) UpdateProvider(data dto.UpdateProviderRequest) (*dto.UpdateProviderResponse, error) {

	config := &domains.LLMProviderConfig{
		ID:              data.ID,
		Compatible:      data.Compatible,
		Name:            data.Name,
		BaseURL:         data.BaseURL,
		APIKey:          data.APIKey,
		IsEnable:        data.IsEnable,
		ModelName:       data.ModelName,
		AvailableModels: data.AvailableModels,
	}

	err := global.GlobalConfigMgr.SaveLLMProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ProviderService) CreateProvider(data dto.CreateProviderRequest) (*dto.CreateProviderResponse, error) {
	config := &domains.LLMProviderConfig{
		ID:              "",
		Compatible:      data.Compatible,
		Name:            data.Name,
		BaseURL:         data.BaseURL,
		APIKey:          data.APIKey,
		IsEnable:        data.IsEnable,
		ModelName:       "",
		AvailableModels: data.AvailableModels,
	}

	err := global.GlobalConfigMgr.CreateLLMProviderConfig(config)
	if err != nil {
		return nil, err
	}

	return &dto.CreateProviderResponse{}, nil
}

func (s *ProviderService) DeleteProvider(id dto.DeleteProviderRequest) (*dto.DeleteProviderResponse, error) {
	err := global.GlobalConfigMgr.DeleteLLMProviderConfig(id.ID)
	if err != nil {
		return nil, err
	}
	return &dto.DeleteProviderResponse{}, nil
}

func (s *ProviderService) FetchModel(data dto.FetchModelRequest) (*dto.FetchModelResponse, error) {
	models, err := global.GlobalLLMManager.GetAvailableModels(data.Compatible, data.BaseURL, data.APIKey)
	if err != nil {
		return nil, err
	}
	return &dto.FetchModelResponse{Models: models}, nil
}

// @Injectable
func NewProviderService() *ProviderService {
	return &ProviderService{}
}

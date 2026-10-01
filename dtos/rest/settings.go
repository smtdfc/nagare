package rest

const (
	GetGeneralSettings = "/api/v1/user/settings/general/get"
	SetGeneralSettings = "/api/v1/user/settings/general/set"
)

type GeneralSettings struct {
	DefaultLLMModel    string `json:"defaultLLMModel"`
	DefaultLLMProvider string `json:"defaultLLMProvider"`
}

type GetGeneralSettingsResponse struct {
	GeneralSettings *GeneralSettings `json:"generalSettings"`
}

type SetGeneralSettingsRequest struct {
	GeneralSettings *GeneralSettings `json:"generalSettings"`
}

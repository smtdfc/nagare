package rest

const (
	GetGeneralSettings = "/api/v1/user/settings/general/get"
	SetGeneralSettings = "/api/v1/user/settings/general/set"
)

type GeneralSettings struct {
	CurrentModel    string `json:"current_model"`
	CurrentProvider string `json:"current_provider"`
}

type GetGeneralSettingsResponse struct {
	GeneralSettings *GeneralSettings `json:"general_settings"`
}

type SetGeneralSettingsRequest struct {
	GeneralSettings *GeneralSettings `json:"general_settings"`
}

package rest

const (
	GetGeneralSettings = "/api/v1/user/settings/general/get"
	SetGeneralSettings = "/api/v1/user/settings/general/set"
)

type GeneralSettings struct {
	DefaultModel         string `json:"defaultModel"`
	DefaultModelProvider string `json:"defaultModelProvider"`
}

type GetGeneralSettingsResponse struct {
	GeneralSettings *GeneralSettings `json:"generalSettings"`
}

type SetGeneralSettingsRequest struct {
	GeneralSettings *GeneralSettings `json:"generalSettings"`
}

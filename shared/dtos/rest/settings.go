package rest

const (
	GetGeneralSettings = "/api/v1/user/settings/general/get"
	SetGeneralSettings = "/api/v1/user/settings/general/set"
)

type GeneralSettings struct {
	CurrentModel    string `json:"currentModel"`
	CurrentProvider string `json:"currentProvider"`
}

type GetGeneralSettingsResponse struct {
	GeneralSettings *GeneralSettings `json:"generalSettings"`
}

type SetGeneralSettingsRequest struct {
	GeneralSettings *GeneralSettings `json:"generalSettings"`
}

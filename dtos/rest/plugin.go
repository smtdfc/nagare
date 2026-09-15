package rest

type Plugin struct {
	ID       string   `json:"id"`
	PluginID string   `json:"pluginID"`
	Name     string   `json:"name"`
	Author   string   `json:"author"`
	Features []string `json:"features"`
	Version  string   `json:"version"`
	IsActive bool     `json:"isActive"`
}

const (
	GetListPluginEndpoint      = "/api/v1/user/plugins/list"
	InstallLocalPluginEndpoint = "/api/v1/user/plugins/install-local"
	UninstallPluginEndpoint    = "/api/v1/user/plugins/uninstall"
	ActivatePluginEndpoint     = "/api/v1/user/plugins/activate"
	DeactivatePluginEndpoint   = "/api/v1/user/plugins/deactivate"
)

type GetListPluginResponse struct {
	Plugins []*Plugin `json:"plugins"`
}

type InstallLocalPluginRequest struct {
	Path string `json:"path"`
}

type UninstallPluginRequest struct {
	ID string `json:"id"`
}

type ActivatePluginRequest struct {
	ID string `json:"id"`
}

type DeactivatePluginRequest struct {
	ID string `json:"id"`
}

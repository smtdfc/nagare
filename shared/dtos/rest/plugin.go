package rest

type Plugin struct {
	ID       string   `json:"id"`
	PluginID string   `json:"plugin_id"`
	Name     string   `json:"name"`
	Author   string   `json:"author"`
	Features []string `json:"features"`
	Version  string   `json:"version"`
	IsActive bool     `json:"is_active"`
}

const (
	GetListPluginEndpoint      = "/api/v1/user/plugins/list"
	InstallLocalPluginEndpoint = "/api/v1/user/plugins/install-local"
)

type GetListPluginResponse struct {
	Plugins []*Plugin `json:"plugins"`
}

type InstallLocalPluginRequest struct {
	Path string `json:"path"`
}

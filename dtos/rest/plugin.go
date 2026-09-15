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

type PluginStatus struct {
	PID         string  `json:"pid"`
	PluginID    string  `json:"pluginID"`
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	CPUPercent  float64 `json:"cpuPercent"`
	MemoryUsage float64 `json:"memoryUsage"`
}

const (
	GetListPluginEndpoint      = "/api/v1/user/plugins/list"
	InstallLocalPluginEndpoint = "/api/v1/user/plugins/install-local"
	UninstallPluginEndpoint    = "/api/v1/user/plugins/uninstall"
	ActivatePluginEndpoint     = "/api/v1/user/plugins/activate"
	DeactivatePluginEndpoint   = "/api/v1/user/plugins/deactivate"
	GetPluginStatusEndpoint    = "/api/v1/user/plugins/status"
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

type GetPluginStatusRequest struct {
	ID string `json:"id"`
}

type GetPluginStatusResponse struct {
	Status *PluginStatus `json:"status"`
}

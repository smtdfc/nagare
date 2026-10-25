package rest

type Plugin struct {
	ID              string   `json:"id"`
	PackageName     string   `json:"packageName"`
	Name            string   `json:"name"`
	Author          string   `json:"author"`
	Features        []string `json:"features"`
	Version         string   `json:"version"`
	SignatureStatus string   `json:"signatureStatus"`
	IsActive        bool     `json:"isActive"`
}

type PluginStatus struct {
	PID         string  `json:"pid"`
	PackageName string  `json:"packageName"`
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	CPUPercent  float64 `json:"cpuPercent"`
	MemoryUsage float64 `json:"memoryUsage"`
}

const (
	GetListPluginEndpoint               = "/api/v1/user/plugins/list"
	InstallLocalPluginEndpoint          = "/api/v1/user/plugins/install-local"
	InstallPluginFromAttachmentEndpoint = "/api/v1/user/plugins/install-attachment"
	UpdatePluginEndpoint                = "/api/v1/user/plugins/update"
	UninstallPluginEndpoint             = "/api/v1/user/plugins/uninstall"
	ActivatePluginEndpoint              = "/api/v1/user/plugins/activate"
	DeactivatePluginEndpoint            = "/api/v1/user/plugins/deactivate"
	GetPluginStatusEndpoint             = "/api/v1/user/plugins/status"
	UploadPluginEndpoint                = "/api/v1/user/plugins/upload"
)

type GetListPluginResponse struct {
	Plugins []*Plugin `json:"plugins"`
}

type InstallLocalPluginRequest struct {
	Path string `json:"path"`
}

type InstallLocalPluginResponse struct {
	Plugin  *Plugin `json:"plugin"`
	Warning string  `json:"warning,omitempty"`
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

type UploadPluginResponse struct {
	AttachmentID string `json:"attachmentId"`
}

type InstallPluginFromAttachmentRequest struct {
	AttachmentID          string `json:"attachmentId"`
	SignatureAttachmentID string `json:"signatureAttachmentId,omitempty"`
}

type InstallPluginFromAttachmentResponse struct {
	Plugin  *Plugin `json:"plugin"`
	Warning string  `json:"warning,omitempty"`
}

type UpdatePluginRequest struct {
	AttachmentID          string `json:"attachmentId"`
	SignatureAttachmentID string `json:"signatureAttachmentId,omitempty"`
}

type UpdatePluginResponse struct {
	Plugin  *Plugin `json:"plugin"`
	Warning string  `json:"warning,omitempty"`
}

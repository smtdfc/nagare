package dto

type PluginHandshakeRequest struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	ApiVersion string `json:"api_version"`
}

type PluginHandshakeResponse struct {
	Accepted bool   `json:"accepted"`
	Message  string `json:"message"`
}

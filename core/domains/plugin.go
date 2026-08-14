package domains

type PluginInfo struct {
	ID         string
	PluginID   string
	Name       string
	Active     bool
	ApiVersion string
	Author     string
	Version    string
	Bin        string
}

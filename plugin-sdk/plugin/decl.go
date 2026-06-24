package plugin

type PluginMetadata struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Author   string   `json:"author"`
	Version  string   `json:"version"`
	Features []string `json:"features"`
	Bin      string   `json:"bin"`
}

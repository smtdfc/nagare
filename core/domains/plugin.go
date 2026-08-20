package domains

import (
	"github.com/smtdfc/nagare/shared/plugin"
)

type PluginInfo struct {
	ID         string
	PluginID   string
	Name       string
	Active     bool
	ApiVersion string
	Author     string
	Version    string
	Bin        string
	Features   plugin.ListPluginFeature
}

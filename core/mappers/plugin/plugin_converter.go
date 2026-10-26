package plugin

import (
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/plugin"
)

func parsePluginFeatures(s string) []plugin.Feature {
	return plugin.ParseFeatureString(s)
}

func formatPluginFeatures(features []plugin.Feature) string {
	p := &plugin.Plugin{Features: features}
	return p.ToFeaturesString()
}

func mapSignatureStatus(status string) string {
	if status == "" {
		return plugin.SignatureStatusUnsigned
	}
	return status
}

func mapPluginStatus(status string) string {
	if status == "" {
		return plugin.PluginStatusUnknown
	}
	return status
}

// goverter:converter
// goverter:output:file ./mapper.generated.go
// goverter:ignoreMissing
type IPluginMapper interface {
	// goverter:map Features | parsePluginFeatures
	// goverter:map SignatureStatus | mapSignatureStatus
	// goverter:map LastStatus | mapPluginStatus
	ToDomain(entity *entities.Plugin) *plugin.Plugin

	// goverter:map Features | formatPluginFeatures
	ToEntity(domain *plugin.Plugin) *entities.Plugin

	ToDomains(entities []*entities.Plugin) []*plugin.Plugin
}

package host

import "github.com/smtdfc/nagare/core/plugin"

type PluginGroup string

var (
	VectorIndexPluginGroup PluginGroup = "VectorIndexPluginGroup"
	NormalPluginGroup      PluginGroup = "NormalPluginGroup"
)

func getGroupName(features []plugin.Feature) PluginGroup {
	for _, f := range features {
		if f == plugin.VectorIndexFeature {
			return VectorIndexPluginGroup
		}
	}

	return NormalPluginGroup
}

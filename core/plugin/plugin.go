package plugin

import "strings"

type PluginFeature string

const (
	CHAT_FEATURE PluginFeature = "CHAT_FEATURE"
)

func (p PluginFeature) ToString() string {
	return string(p)
}

func ParseFeatureString(raw string) []PluginFeature {
	parts := strings.Split(raw, ",")
	features := []PluginFeature{}
	for _, p := range parts {
		switch p {
		case string(CHAT_FEATURE):
			features = append(features, CHAT_FEATURE)
		}
	}

	return features
}

type Plugin struct {
	ID       string
	PluginID string
	Name     string
	Author   string
	Features []PluginFeature
	Version  string
	Bin      string
	IsActive bool
}

func (p *Plugin) ToFeaturesString() string {
	var s strings.Builder
	for _, f := range p.Features {
		s.WriteString(f.ToString())
	}

	return s.String()
}

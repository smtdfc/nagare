package plugin

import (
	"strings"

	"github.com/google/uuid"
)

type Feature string

const (
	ChatFeature Feature = "CHAT_FEATURE"
)

func (p Feature) ToString() string {
	return string(p)
}

func ParseFeatureString(raw string) []Feature {
	parts := strings.Split(raw, ",")
	var features []Feature
	for _, p := range parts {
		switch p {
		case string(ChatFeature):
			features = append(features, ChatFeature)
		}
	}

	return features
}

type Plugin struct {
	ID       uuid.UUID
	PluginID string
	Name     string
	Author   string
	Features []Feature
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

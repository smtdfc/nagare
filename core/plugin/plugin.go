package plugin

import (
	"strings"

	"github.com/google/uuid"
)

type Feature string

const (
	ChatFeature        Feature = "chat"
	DynamicToolFeature Feature = "dynamic_tool"
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
		case string(DynamicToolFeature):
			features = append(features, DynamicToolFeature)
		}
	}

	return features
}

type Plugin struct {
	ID          uuid.UUID
	PackageName string
	Name        string
	Author      string
	Features    []Feature
	Version     string
	Bin         string
	IsActive    bool
}

func (p *Plugin) ToFeaturesString() string {
	var s strings.Builder
	for _, f := range p.Features {
		s.WriteString(f.ToString())
	}

	return s.String()
}

type PluginStatus struct {
	PID         string
	PackageName string
	Name        string
	Version     string
	CPUPercent  float64
	MemoryUsage float64
}

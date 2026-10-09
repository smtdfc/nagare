package plugin

import (
	"strings"

	"github.com/google/uuid"
)

type Feature string

const (
	ChatFeature Feature = "chat"
	ToolFeature Feature = "plugin_tool"
)

func (p Feature) ToString() string {
	return string(p)
}

func ParseFeatureString(raw string) []Feature {
	parts := strings.Split(raw, ",")
	var features []Feature
	for _, p := range parts {
		p = strings.TrimSpace(p)
		switch p {
		case string(ChatFeature):
			features = append(features, ChatFeature)
		case string(ToolFeature):
			features = append(features, ToolFeature)
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
	strs := make([]string, len(p.Features))
	for i, f := range p.Features {
		strs[i] = f.ToString()
	}
	return strings.Join(strs, ",")
}

type Status struct {
	PID         string
	PackageName string
	Name        string
	Version     string
	CPUPercent  float64
	MemoryUsage float64
}

package plugin

import (
	"strings"

	"github.com/google/uuid"
)

type Feature string

const (
	ChatFeature          Feature = "chat"
	ModelProviderFeature Feature = "model_provider"
	ToolFeature          Feature = "plugin_tool"
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
		case string(ModelProviderFeature):
			features = append(features, ModelProviderFeature)
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
	BinChecksum string
	// SignatureStatus records how the installed archive was verified:
	// "verified", "unverified", or "unsigned".
	SignatureStatus string
	IsActive        bool
}

const (
	SignatureStatusVerified   = "verified"
	SignatureStatusUnverified = "unverified"
	SignatureStatusUnsigned   = "unsigned"
)

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

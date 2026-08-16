package plugin

import "errors"

type PluginBin struct {
	Architecture string `json:"architecture"`
	Path         string `json:"path"`
}

type PluginMetadata struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Version       string      `json:"version"`
	ApiVersion    string      `json:"api_version"`
	Author        string      `json:"author"`
	Architectures []string    `json:"architectures"`
	Bins          []PluginBin `json:"bins"`
}

func (m *PluginMetadata) SupportsArchitecture(arch string) bool {
	for _, bin := range m.Bins {
		if bin.Architecture == arch {
			return true
		}
	}
	return false
}

func (m *PluginMetadata) GetBinForArchitecture(arch string) (*PluginBin, error) {
	for i := range m.Bins {
		if m.Bins[i].Architecture == arch {
			return &m.Bins[i], nil
		}
	}
	return nil, errors.New("binary for architecture not found")
}

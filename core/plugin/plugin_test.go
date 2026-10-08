package plugin

import (
	"testing"

	"github.com/google/uuid"
)

func TestFeature_ToString(t *testing.T) {
	// Verify Feature ToString method
	if ChatFeature.ToString() != "chat" {
		t.Errorf("expected 'chat', got '%s'", ChatFeature.ToString())
	}
	if PluginToolFeature.ToString() != "plugin_tool" {
		t.Errorf("expected 'plugin_tool', got '%s'", PluginToolFeature.ToString())
	}
}

func TestParseFeatureString(t *testing.T) {
	// Verify parsing features from comma-separated string
	tests := []struct {
		input    string
		expected []Feature
	}{
		{"chat,plugin_tool", []Feature{ChatFeature, PluginToolFeature}},
		{"chat", []Feature{ChatFeature}},
		{"plugin_tool", []Feature{PluginToolFeature}},
		{" chat , plugin_tool ", []Feature{ChatFeature, PluginToolFeature}},
		{"unknown,chat", []Feature{ChatFeature}},
		{"", nil},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			res := ParseFeatureString(tc.input)
			if len(res) != len(tc.expected) {
				t.Fatalf("expected %d features, got %d", len(tc.expected), len(res))
			}
			for i := range res {
				if res[i] != tc.expected[i] {
					t.Errorf("expected feature %s, got %s", tc.expected[i], res[i])
				}
			}
		})
	}
}

func TestPlugin_ToFeaturesString(t *testing.T) {
	// Verify Plugin.ToFeaturesString method
	id := uuid.New()
	p := Plugin{
		ID:          id,
		PackageName: "com.example.telegram",
		Name:        "Telegram",
		Author:      "Author",
		Features:    []Feature{ChatFeature, PluginToolFeature},
		Version:     "1.0.0",
		Bin:         "telegram.bin",
		IsActive:    true,
	}

	res := p.ToFeaturesString()
	if res != "chat,plugin_tool" {
		t.Errorf("expected 'chat,plugin_tool', got '%s'", res)
	}
}

func TestPluginStatus(t *testing.T) {
	// Verify PluginStatus struct fields
	status := PluginStatus{
		PID:         "12345",
		PackageName: "com.example.test",
		Name:        "Test",
		Version:     "0.1.0",
		CPUPercent:  1.5,
		MemoryUsage: 25.4,
	}

	if status.PID != "12345" {
		t.Errorf("expected PID '12345', got '%s'", status.PID)
	}
	if status.CPUPercent != 1.5 {
		t.Errorf("expected CPU 1.5, got %f", status.CPUPercent)
	}
}

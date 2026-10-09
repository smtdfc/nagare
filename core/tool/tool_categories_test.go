package tool

import (
	"testing"
)

func TestToolCategoryConstants(t *testing.T) {
	// Verify tool category constants values
	tests := []struct {
		name     string
		category string
		expected string
	}{
		{"ProcessManagement", ProcessManagementCategory, "process_management"},
		{"PowerManagement", PowerManagementCategory, "power_control"},
		{"AudioManagement", AudioManagementCategory, "audio_management"},
		{"Filesystem", FilesystemCategory, "filesystem"},
		{"Browser", BrowserCategory, "browser"},
		{"Weather", WeatherCategory, "weather"},
		{"Networking", NetworkingCategory, "networking"},
		{"TaskManagement", TaskManagementCategory, "task_management"},
		{"ToolRouting", RoutingCategory, "tool_routing"},
		{"Timing", TimingCategory, "timing"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.category != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, tc.category)
			}
		})
	}
}

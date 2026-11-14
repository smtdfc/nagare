package manager

import (
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestGetCategoriesPrompt_MergesPluginCategories(t *testing.T) {
	mgr := &ToolManager{
		toolCategories: map[string]string{
			tool.FilesystemCategory: "Filesystem",
			tool.TimingCategory:     "Timing",
		},
		pluginToolCategories: map[string]map[string]string{
			"plugin-a": {
				"plugin_echo":       "Plugin Echo",
				tool.TimingCategory: "Hijacked Timing",
			},
		},
	}

	prompt := mgr.GetCategoriesPrompt()

	if !strings.Contains(prompt, "plugin_echo (Plugin Echo)") {
		t.Errorf("expected plugin category in prompt, got %q", prompt)
	}
	if strings.Contains(prompt, "Hijacked Timing") {
		t.Errorf("expected built-in category description to win, got %q", prompt)
	}

	echoIndex := strings.Index(prompt, "plugin_echo")
	filesystemIndex := strings.Index(prompt, tool.FilesystemCategory)
	timingIndex := strings.Index(prompt, tool.TimingCategory)
	if filesystemIndex > echoIndex || echoIndex > timingIndex {
		t.Errorf("expected sorted category prompt, got %q", prompt)
	}
}

func TestRemovePluginToolCategories(t *testing.T) {
	mgr := &ToolManager{
		toolCategories: map[string]string{},
		pluginToolCategories: map[string]map[string]string{
			"plugin-a": {"plugin_echo": "Plugin Echo"},
		},
	}

	mgr.RemovePluginToolCategories("plugin-a")

	if strings.Contains(mgr.GetCategoriesPrompt(), "plugin_echo") {
		t.Error("expected removed plugin categories to disappear from the prompt")
	}
}

package prompt

import (
	"strings"
	"testing"
)

func TestSystemPromptTemplate_Markers(t *testing.T) {
	// Verify systemPromptContent contains expected structure and tags
	if len(systemPromptContent) == 0 {
		t.Fatalf("expected non-empty systemPromptContent")
	}

	expectedMarkers := []string{
		"You are Nagare",
		"Mandatory Rules",
		"Tool Selection and Execution",
		"find_tool_by_categories",
		"execute_tool",
		"{{.ToolCategories}}",
		"Accurate and Natural Responses",
		"Mathematical Expressions and UI Rendering",
		"Additional Rules",
		"{{.PluginPrompts}}",
	}

	for _, marker := range expectedMarkers {
		if !strings.Contains(systemPromptContent, marker) {
			t.Errorf("expected prompt content to contain '%s'", marker)
		}
	}
}

func TestSystemPromptTemplate_Build(t *testing.T) {
	// Verify SystemPromptTemplate rendering with ToolCategories and PluginPrompts
	categories := "- filesystem\n- browser\n- weather"
	pluginPrompts := "0: extra context \n"
	rendered, err := SystemPromptTemplate.Build(map[string]any{
		"ToolCategories": categories,
		"PluginPrompts":  pluginPrompts,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(rendered, categories) {
		t.Errorf("expected rendered prompt to contain tool categories")
	}
	if !strings.Contains(rendered, "Nagare") {
		t.Errorf("expected rendered prompt to contain 'Nagare'")
	}
	if !strings.Contains(rendered, pluginPrompts) {
		t.Errorf("expected rendered prompt to contain plugin prompts")
	}
}

func TestBuildSystemPrompt_WithPluginPrompt(t *testing.T) {
	rendered, err := BuildSystemPrompt("- filesystem", []PluginSystemPrompt{
		{PluginID: "plugin-1", Content: "Current chat is a group named Release."},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(rendered, "- filesystem") {
		t.Errorf("expected rendered prompt to contain tool categories")
	}
	if !strings.Contains(rendered, "0: Current chat is a group named Release. \n") {
		t.Errorf("expected rendered prompt to contain numbered plugin prompt")
	}
}

func TestBuildSystemPrompt_WithMultiplePluginPrompts(t *testing.T) {
	rendered, err := BuildSystemPrompt("- filesystem", []PluginSystemPrompt{
		{PluginID: "plugin-1", Content: "first context"},
		{PluginID: "plugin-2", Content: "second context"},
		{PluginID: "plugin-3", Content: "third context"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantOrder := []string{
		"0: first context \n",
		"1: second context \n",
		"2: third context \n",
	}

	lastIndex := -1
	for _, want := range wantOrder {
		index := strings.Index(rendered, want)
		if index == -1 {
			t.Fatalf("expected rendered prompt to contain %q", want)
		}
		if index <= lastIndex {
			t.Errorf("expected %q to appear after the previous plugin prompt", want)
		}
		lastIndex = index
	}
}

func TestBuildSystemPrompt_WithoutPluginPrompt(t *testing.T) {
	rendered, err := BuildSystemPrompt("- filesystem", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasSuffix(strings.TrimSpace(rendered), "## Additional Rules") {
		t.Errorf("expected empty plugin prompts to leave the additional rules section empty")
	}
}

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
		"<system_instructions>",
		"</system_instructions>",
		"Nagare",
		"tool_routing",
		"find_tool_by_categories",
		"execute_tool",
		"{{.ToolCategories}}",
		"response_language",
	}

	for _, marker := range expectedMarkers {
		if !strings.Contains(systemPromptContent, marker) {
			t.Errorf("expected prompt content to contain '%s'", marker)
		}
	}
}

func TestSystemPromptTemplate_Build(t *testing.T) {
	// Verify SystemPromptTemplate rendering with ToolCategories
	categories := "- filesystem\n- browser\n- weather"
	rendered, err := SystemPromptTemplate.Build(map[string]any{
		"ToolCategories": categories,
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
}

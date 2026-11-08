package prompt

import "testing"

func TestPluginPromptStore_SetAndList(t *testing.T) {
	store := NewPluginPromptStore()

	store.Set("plugin-2", "second context")
	store.Set("plugin-1", "first context")

	prompts := store.List()
	if len(prompts) != 2 {
		t.Fatalf("expected 2 prompts, got %d", len(prompts))
	}

	if prompts[0].PluginID != "plugin-1" || prompts[0].Content != "first context" {
		t.Fatalf("expected plugin-1 first, got %+v", prompts[0])
	}
	if prompts[1].PluginID != "plugin-2" || prompts[1].Content != "second context" {
		t.Fatalf("expected plugin-2 second, got %+v", prompts[1])
	}

	store.Set("plugin-1", "replaced context")
	prompts = store.List()
	if len(prompts) != 2 {
		t.Fatalf("expected 2 prompts after replace, got %d", len(prompts))
	}
	if prompts[0].Content != "replaced context" {
		t.Fatalf("expected replaced content, got %q", prompts[0].Content)
	}
}

func TestPluginPromptStore_Clear(t *testing.T) {
	store := NewPluginPromptStore()
	store.Set("plugin-1", "plugin one")
	store.Set("plugin-2", "plugin two")

	store.Set("plugin-1", "   ")
	prompts := store.List()
	if len(prompts) != 1 || prompts[0].PluginID != "plugin-2" {
		t.Fatalf("whitespace content should clear the fragment, got %+v", prompts)
	}

	store.Set("plugin-1", "plugin one again")
	store.Clear("plugin-1")
	prompts = store.List()
	if len(prompts) != 1 || prompts[0].PluginID != "plugin-2" {
		t.Fatalf("expected only plugin-2 after clear, got %+v", prompts)
	}

	store.ClearAll()
	if prompts := store.List(); len(prompts) != 0 {
		t.Fatalf("expected all prompts cleared, got %+v", prompts)
	}
}

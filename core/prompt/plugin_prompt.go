package prompt

import (
	"sort"
	"strings"
	"sync"
)

type PluginSystemPrompt struct {
	PluginID string
	Content  string
}

type PluginPromptStore struct {
	mu      sync.RWMutex
	prompts map[string]string
}

func NewPluginPromptStore() *PluginPromptStore {
	return &PluginPromptStore{
		prompts: make(map[string]string),
	}
}

var defaultPluginPromptStore = NewPluginPromptStore()

func SetPluginSystemPrompt(pluginID, content string) {
	defaultPluginPromptStore.Set(pluginID, content)
}

func GetPluginSystemPrompts() []PluginSystemPrompt {
	return defaultPluginPromptStore.List()
}

func ClearPluginSystemPrompt(pluginID string) {
	defaultPluginPromptStore.Clear(pluginID)
}

func ClearAllPluginSystemPrompts() {
	defaultPluginPromptStore.ClearAll()
}

func (s *PluginPromptStore) Set(pluginID, content string) {
	if strings.TrimSpace(content) == "" {
		s.Clear(pluginID)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.prompts[pluginID] = content
}

func (s *PluginPromptStore) List() []PluginSystemPrompt {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]PluginSystemPrompt, 0, len(s.prompts))
	for pluginID, content := range s.prompts {
		result = append(result, PluginSystemPrompt{
			PluginID: pluginID,
			Content:  content,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].PluginID < result[j].PluginID
	})

	return result
}

func (s *PluginPromptStore) Clear(pluginID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.prompts, pluginID)
}

func (s *PluginPromptStore) ClearAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.prompts = make(map[string]string)
}

package provider

import (
	"testing"
)

type dummyEmbeddingAdapter struct{}

func TestEmbeddingProviderAdapter_Interface(t *testing.T) {
	// Verify mock satisfies EmbeddingProviderAdapter interface
	var adapter EmbeddingProviderAdapter = &dummyEmbeddingAdapter{}
	if adapter == nil {
		t.Errorf("expected non-nil adapter")
	}
}

package manager

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"slices"
	"testing"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/discovery"
	"github.com/smtdfc/nagare/core/logger"
)

type fakeSourceProvider struct {
	searchResult []discovery.Source
	searchErr    error
	fetchResult  *discovery.Detail
	fetchErr     error

	gotQuery string
	gotLimit int
	gotID    string
}

func (f *fakeSourceProvider) Search(_ context.Context, query string, limit int) ([]discovery.Source, error) {
	f.gotQuery = query
	f.gotLimit = limit
	return f.searchResult, f.searchErr
}

func (f *fakeSourceProvider) Fetch(_ context.Context, id string) (*discovery.Detail, error) {
	f.gotID = id
	return f.fetchResult, f.fetchErr
}

func newTestDiscoveryManager() *DiscoveryManager {
	return NewDiscoveryManager(&logger.BaseLogger{
		Logger: *slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
}

func TestDiscoveryManager_DefaultProviders(t *testing.T) {
	mgr := newTestDiscoveryManager()

	if providers := mgr.GetAvailableSourceProvider(); !slices.Contains(providers, "Wikipedia") {
		t.Fatalf("expected Wikipedia provider, got %v", providers)
	}
}

func TestDiscoveryManager_AddProvider(t *testing.T) {
	mgr := newTestDiscoveryManager()
	mgr.AddProvider("Fake", &fakeSourceProvider{})

	if providers := mgr.GetAvailableSourceProvider(); !slices.Contains(providers, "Fake") {
		t.Fatalf("expected Fake provider, got %v", providers)
	}
}

func TestDiscoveryManager_Search(t *testing.T) {
	provider := &fakeSourceProvider{
		searchResult: []discovery.Source{
			{Provider: "Fake", ID: "1", Title: "Golang"},
		},
	}
	mgr := newTestDiscoveryManager()
	mgr.AddProvider("Fake", provider)

	result, err := mgr.Search(context.Background(), "Fake", "golang", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.gotQuery != "golang" {
		t.Errorf("expected query 'golang', got %q", provider.gotQuery)
	}
	if provider.gotLimit != 5 {
		t.Errorf("expected limit 5, got %d", provider.gotLimit)
	}
	if len(result) != 1 || result[0].Title != "Golang" {
		t.Fatalf("expected one Golang result, got %+v", result)
	}
}

func TestDiscoveryManager_Fetch(t *testing.T) {
	provider := &fakeSourceProvider{
		fetchResult: &discovery.Detail{SourceID: "1", Content: "content"},
	}
	mgr := newTestDiscoveryManager()
	mgr.AddProvider("Fake", provider)

	result, err := mgr.Fetch(context.Background(), "Fake", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.gotID != "1" {
		t.Errorf("expected id '1', got %q", provider.gotID)
	}
	if result == nil || result.SourceID != "1" || result.Content != "content" {
		t.Fatalf("expected fetched detail, got %+v", result)
	}
}

func TestDiscoveryManager_UnknownProvider(t *testing.T) {
	mgr := newTestDiscoveryManager()

	if _, err := mgr.Search(context.Background(), "Missing", "query", 1); !errors.Is(err, custom_errors.ErrSourceProviderNotFound) {
		t.Errorf("expected ErrSourceProviderNotFound for search, got %v", err)
	}
	if _, err := mgr.Fetch(context.Background(), "Missing", "1"); !errors.Is(err, custom_errors.ErrSourceProviderNotFound) {
		t.Errorf("expected ErrSourceProviderNotFound for fetch, got %v", err)
	}
}

func TestDiscoveryManager_ProviderFailure(t *testing.T) {
	mgr := newTestDiscoveryManager()
	mgr.AddProvider("Fake", &fakeSourceProvider{
		searchErr: errors.New("search failed"),
		fetchErr:  errors.New("fetch failed"),
	})

	if _, err := mgr.Search(context.Background(), "Fake", "query", 1); !errors.Is(err, custom_errors.ErrDiscoveryFailed) {
		t.Errorf("expected ErrDiscoveryFailed for search, got %v", err)
	}
	if _, err := mgr.Fetch(context.Background(), "Fake", "1"); !errors.Is(err, custom_errors.ErrDiscoveryFailed) {
		t.Errorf("expected ErrDiscoveryFailed for fetch, got %v", err)
	}
}

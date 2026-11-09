package declarations

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/computer_use"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/discovery"
	discovery_manager "github.com/smtdfc/nagare/core/discovery/manager"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool"
)

type fakeDiscoveryProvider struct {
	searchResult []discovery.Source
	searchErr    error
	fetchResult  *discovery.Detail
	fetchErr     error

	gotQuery string
	gotLimit int
	gotID    string
}

func (f *fakeDiscoveryProvider) Search(_ context.Context, query string, limit int) ([]discovery.Source, error) {
	f.gotQuery = query
	f.gotLimit = limit
	return f.searchResult, f.searchErr
}

func (f *fakeDiscoveryProvider) Fetch(_ context.Context, id string) (*discovery.Detail, error) {
	f.gotID = id
	return f.fetchResult, f.fetchErr
}

type fakeDiscoveryBindings struct {
	mockCronJobBindings
	discoveryMgr *discovery_manager.DiscoveryManager
}

// GetComputerUseOrchestrator implements [tool.Bindings].
func (f *fakeDiscoveryBindings) GetComputerUseOrchestrator() *computer_use.ComputerUseOrchestrator {
	panic("unimplemented")
}

func (f *fakeDiscoveryBindings) GetDiscoveryMgr() *discovery_manager.DiscoveryManager {
	return f.discoveryMgr
}

func newFakeDiscoveryBindings(provider discovery.SourceProviderAdapter) *fakeDiscoveryBindings {
	mgr := discovery_manager.NewDiscoveryManager(&logger.BaseLogger{
		Logger: *slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	mgr.AddProvider("Fake", provider)

	return &fakeDiscoveryBindings{discoveryMgr: mgr}
}

func TestSearchTool_Metadata(t *testing.T) {
	if SearchTool.GetName() != "search_tool" {
		t.Errorf("expected 'search_tool', got '%s'", SearchTool.GetName())
	}
	if len(SearchTool.GetCategories()) != 1 || SearchTool.GetCategories()[0] != tool.ResearchCategory {
		t.Errorf("expected ResearchCategory, got %v", SearchTool.GetCategories())
	}

	schema := SearchTool.GetArgsSchema()
	for _, field := range []string{`"provider"`, `"query"`, `"limit"`} {
		if !strings.Contains(schema, field) {
			t.Errorf("expected schema to contain %s, got %s", field, schema)
		}
	}
}

func TestSearchTool_Execute(t *testing.T) {
	provider := &fakeDiscoveryProvider{
		searchResult: []discovery.Source{
			{Provider: "Fake", ID: "1", Title: "Golang", Description: "A programming language"},
		},
	}
	tl := withBindings(t, SearchTool, newFakeDiscoveryBindings(provider))

	res, err := tl.Execute(nil, `{"provider":"Fake","query":"golang","limit":5}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.gotQuery != "golang" {
		t.Errorf("expected query 'golang', got %q", provider.gotQuery)
	}
	if provider.gotLimit != 5 {
		t.Errorf("expected limit 5, got %d", provider.gotLimit)
	}

	var out SearchOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if len(out.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out.Results))
	}
	if out.Results[0].Title != "Golang" {
		t.Errorf("expected title 'Golang', got %q", out.Results[0].Title)
	}
}

func TestSearchTool_PropagatesError(t *testing.T) {
	tl := withBindings(t, SearchTool, newFakeDiscoveryBindings(&fakeDiscoveryProvider{
		searchErr: errors.New("search failed"),
	}))

	_, err := tl.Execute(nil, `{"provider":"Fake","query":"golang","limit":5}`)
	if !errors.Is(err, custom_errors.ErrDiscoveryFailed) {
		t.Fatalf("expected ErrDiscoveryFailed, got %v", err)
	}
}

var _ tool.Bindings = (*fakeDiscoveryBindings)(nil)

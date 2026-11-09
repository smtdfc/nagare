package declarations

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/discovery"
	"github.com/smtdfc/nagare/core/tool"
)

func TestFetchContentTool_Metadata(t *testing.T) {
	if FetchContentTool.GetName() != "fetch_content_tool" {
		t.Errorf("expected 'fetch_content_tool', got '%s'", FetchContentTool.GetName())
	}
	if len(FetchContentTool.GetCategories()) != 1 || FetchContentTool.GetCategories()[0] != tool.ResearchCategory {
		t.Errorf("expected ResearchCategory, got %v", FetchContentTool.GetCategories())
	}

	schema := FetchContentTool.GetArgsSchema()
	for _, field := range []string{`"provider"`, `"id"`} {
		if !strings.Contains(schema, field) {
			t.Errorf("expected schema to contain %s, got %s", field, schema)
		}
	}
}

func TestFetchContentTool_Execute(t *testing.T) {
	provider := &fakeDiscoveryProvider{
		fetchResult: &discovery.Detail{SourceID: "1", Content: "full content"},
	}
	tl := withBindings(t, FetchContentTool, newFakeDiscoveryBindings(provider))

	res, err := tl.Execute(nil, `{"provider":"Fake","id":"1"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.gotID != "1" {
		t.Errorf("expected id '1', got %q", provider.gotID)
	}

	var out FetchContentOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	if out.Result == nil || out.Result.SourceID != "1" || out.Result.Content != "full content" {
		t.Fatalf("expected fetched content, got %+v", out.Result)
	}
}

func TestFetchContentTool_PropagatesError(t *testing.T) {
	tl := withBindings(t, FetchContentTool, newFakeDiscoveryBindings(&fakeDiscoveryProvider{
		fetchErr: errors.New("fetch failed"),
	}))

	_, err := tl.Execute(nil, `{"provider":"Fake","id":"1"}`)
	if !errors.Is(err, custom_errors.ErrDiscoveryFailed) {
		t.Fatalf("expected ErrDiscoveryFailed, got %v", err)
	}
}

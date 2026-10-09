package client

import (
	"context"
	"errors"
	"testing"
)

type sampleToolArgs struct {
	Query string `json:"query"`
}

type sampleToolResult struct {
	Answer string `json:"answer"`
}

func TestDefineToolAndExecution(t *testing.T) {
	tool := DefineTool("search_tool", "searches things", func(ctx *context.Context, args *sampleToolArgs) (*sampleToolResult, error) {
		return &sampleToolResult{Answer: "result for " + args.Query}, nil
	}, []string{"search"})

	if tool.GetName() != "search_tool" {
		t.Fatalf("expected name 'search_tool', got %q", tool.GetName())
	}
	if tool.GetDescription() != "searches things" {
		t.Fatalf("expected description 'searches things', got %q", tool.GetDescription())
	}
	cats := tool.GetCategories()
	if len(cats) != 1 || cats[0] != "search" {
		t.Fatalf("unexpected categories: %+v", cats)
	}

	// Test schema reflection and caching
	schema1 := tool.GetArgsSchema()
	if schema1 == "" || schema1 == "{}" {
		t.Fatalf("expected non-empty json schema, got %q", schema1)
	}
	schema2 := tool.GetArgsSchema()
	if schema1 != schema2 {
		t.Fatalf("schema cache mismatch: %q vs %q", schema1, schema2)
	}

	// Test successful tool execution
	ctx := context.Background()
	res, err := tool.Execute(&ctx, `{"query":"golang"}`)
	if err != nil {
		t.Fatalf("unexpected tool execution error: %v", err)
	}
	expected := `{"answer":"result for golang"}`
	if res != expected {
		t.Fatalf("expected result %q, got %q", expected, res)
	}

	// Test execution with invalid json arguments
	_, err = tool.Execute(&ctx, `{invalid_json`)
	if !errors.Is(err, ErrIncorrectToolArgs) {
		t.Fatalf("expected ErrIncorrectToolArgs on malformed json, got %v", err)
	}
}

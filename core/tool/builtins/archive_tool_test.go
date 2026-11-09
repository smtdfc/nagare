package declarations

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveTool_Metadata(t *testing.T) {
	if ArchiveTool.GetName() != "archive_tool" {
		t.Errorf("expected 'archive_tool', got '%s'", ArchiveTool.GetName())
	}
	if len(ArchiveTool.GetCategories()) == 0 || ArchiveTool.GetCategories()[0] != "filesystem" {
		t.Errorf("expected FilesystemCategory, got %v", ArchiveTool.GetCategories())
	}
}

func TestArchiveTool_CreateAndExtractZip(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, ArchiveTool, approvals)

	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "inner", "note.txt")
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, []byte("archive me"), 0o644); err != nil {
		t.Fatal(err)
	}

	archivePath := filepath.Join(tmp, "bundle.zip")
	res, err := ArchiveTool.Execute(ctx, `{"action":"create","archive_path":"`+archivePath+`","sources":["`+filepath.Join(tmp, "inner")+`"],"confirm":true}`)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}
	if !contains(res, "zip") {
		t.Errorf("expected format zip in %s", res)
	}
	if _, err := os.Stat(archivePath); err != nil {
		t.Errorf("expected archive to exist: %v", err)
	}

	extractDir := filepath.Join(t.TempDir(), "out")
	res, err = ArchiveTool.Execute(ctx, `{"action":"extract","archive_path":"`+archivePath+`","destination":"`+extractDir+`","confirm":true}`)
	if err != nil {
		t.Fatalf("unexpected extract error: %v", err)
	}
	if len(approvals.calls) != 2 {
		t.Errorf("expected two approval requests, got %d", len(approvals.calls))
	}

	extracted := filepath.Join(extractDir, "inner", "note.txt")
	data, err := os.ReadFile(extracted)
	if err != nil || string(data) != "archive me" {
		t.Errorf("expected restored file at %s, got '%s' (err: %v)", extracted, string(data), err)
	}
}

func TestArchiveTool_ExtractZipSlip(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, ArchiveTool, approvals)

	// Craft an archive whose entry escapes the destination.
	evil := filepath.Join(t.TempDir(), "evil.zip")
	f, err := os.Create(evil)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	h := &zip.FileHeader{Name: "../escaped.txt"}
	zw.CreateHeader(h)
	zw.CreateHeader(&zip.FileHeader{Name: "ok.txt"})
	zw.Close()
	f.Close()

	_, err = ArchiveTool.Execute(ctx, `{"action":"extract","archive_path":"`+evil+`","confirm":true}`)
	if err == nil {
		t.Fatalf("expected zip-slip extraction to be refused")
	}
	if !contains(err.Error(), "zip-slip") {
		t.Errorf("expected a zip-slip error, got: %v", err)
	}
}

func TestArchiveTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	withBindings(t, ArchiveTool, &fakeApprovalBindings{approve: true})

	cases := []struct {
		name string
		args string
	}{
		{"missing confirm", `{"action":"create","archive_path":"/tmp/a.zip"}`},
		{"unknown action", `{"action":"append","archive_path":"/tmp/a.zip","confirm":true}`},
		{"missing archive path", `{"action":"create","confirm":true}`},
		{"bad format", `{"action":"create","archive_path":"/tmp/a.rar","sources":["/tmp"],"confirm":true}`},
		{"no sources", `{"action":"create","archive_path":"/tmp/a.zip","confirm":true}`},
	}
	for _, tc := range cases {
		if _, err := ArchiveTool.Execute(ctx, tc.args); err == nil {
			t.Errorf("expected error for case '%s'", tc.name)
		}
	}
}

func TestArchiveTool_Denied(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: false}
	withBindings(t, ArchiveTool, approvals)

	archivePath := filepath.Join(t.TempDir(), "denied.zip")
	if _, err := ArchiveTool.Execute(ctx, `{"action":"create","archive_path":"`+archivePath+`","sources":["`+t.TempDir()+`"],"confirm":true}`); err == nil {
		t.Fatalf("expected error when the user denies the action")
	}
	if _, ferr := os.Stat(archivePath); !os.IsNotExist(ferr) {
		t.Errorf("expected no archive to be created after denial, got err: %v", ferr)
	}
}

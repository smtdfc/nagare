package upload

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/logger"
)

func TestUploadManager_Lifecycle(t *testing.T) {
	// Verify UploadManager AddAttachment, GetAttachmentPath, and RemoveAttachment
	baseLogger := &logger.BaseLogger{Logger: *slog.Default()}
	um := NewUploadManager(baseLogger)

	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "source.txt")
	err := os.WriteFile(sourceFile, []byte("attachment content"), 0644)
	if err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	id, err := um.AddAttachment(sourceFile)
	if err != nil {
		t.Fatalf("unexpected error adding attachment: %v", err)
	}
	if id == "" {
		t.Fatalf("expected non-empty attachment ID")
	}

	uploadPath, exists := um.GetAttachmentPath(id)
	if !exists {
		t.Errorf("expected attachment to exist")
	}
	if uploadPath == "" {
		t.Errorf("expected valid upload path")
	}

	// Remove attachment
	um.RemoveAttachment(id)
	_, exists = um.GetAttachmentPath(id)
	if exists {
		t.Errorf("expected attachment to be removed")
	}
}

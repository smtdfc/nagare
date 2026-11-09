package declarations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestStatFileTool_Metadata(t *testing.T) {
	assert.Equal(t, "stat_file_tool", StatFileTool.GetName())
	assert.Contains(t, StatFileTool.GetCategories(), tool.FilesystemCategory)
}

func TestStatFileTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	_, err := StatFileTool.Execute(ctx, `{"path":""}`)
	assert.Error(t, err)
}

func TestStatFileTool_NonExistent(t *testing.T) {
	ctx := approvalTestContext()
	rawOut, err := StatFileTool.Execute(ctx, `{"path":"/path/that/does/not/exist_12345"}`)
	assert.NoError(t, err)

	var out StatFileOutput
	err = json.Unmarshal([]byte(rawOut), &out)
	assert.NoError(t, err)
	assert.False(t, out.Exists)
}

func TestStatFileTool_FileAndDir(t *testing.T) {
	ctx := approvalTestContext()
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "info.txt")

	err := os.WriteFile(testFile, []byte("some bytes"), 0644)
	assert.NoError(t, err)

	// Stat File
	rawOut, err := StatFileTool.Execute(ctx, `{"path":"`+testFile+`"}`)
	assert.NoError(t, err)
	var fileOut StatFileOutput
	err = json.Unmarshal([]byte(rawOut), &fileOut)
	assert.NoError(t, err)
	assert.True(t, fileOut.Exists)
	assert.False(t, fileOut.IsDir)
	assert.Equal(t, int64(len("some bytes")), fileOut.Size)
	assert.Equal(t, "info.txt", fileOut.Name)

	// Stat Directory
	rawOutDir, err := StatFileTool.Execute(ctx, `{"path":"`+tmpDir+`"}`)
	assert.NoError(t, err)
	var dirOut StatFileOutput
	err = json.Unmarshal([]byte(rawOutDir), &dirOut)
	assert.NoError(t, err)
	assert.True(t, dirOut.Exists)
	assert.True(t, dirOut.IsDir)
}

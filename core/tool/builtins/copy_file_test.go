package declarations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestCopyFileTool_Metadata(t *testing.T) {
	assert.Equal(t, "copy_file_tool", CopyFileTool.GetName())
	assert.Contains(t, CopyFileTool.GetCategories(), tool.FilesystemCategory)
}

func TestCopyFileTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, CopyFileTool, approvals)

	// confirm = false
	_, err := CopyFileTool.Execute(ctx, `{"source":"a","destination":"b","confirm":false}`)
	assert.Error(t, err)

	// empty source or dest
	_, err = CopyFileTool.Execute(ctx, `{"source":"","destination":"b","confirm":true}`)
	assert.Error(t, err)
}

func TestCopyFileTool_Success(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, CopyFileTool, approvals)

	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "src.txt")
	dstFile := filepath.Join(tmpDir, "sub", "dst.txt")

	err := os.WriteFile(srcFile, []byte("copy this data"), 0644)
	assert.NoError(t, err)

	rawOut, err := CopyFileTool.Execute(ctx, `{"source":"`+srcFile+`","destination":"`+dstFile+`","confirm":true}`)
	assert.NoError(t, err)

	var out CopyFileOutput
	err = json.Unmarshal([]byte(rawOut), &out)
	assert.NoError(t, err)
	assert.Equal(t, int64(len("copy this data")), out.Bytes)

	// Verify content
	copied, err := os.ReadFile(dstFile)
	assert.NoError(t, err)
	assert.Equal(t, "copy this data", string(copied))

	// Attempt overwrite without overwrite=true
	_, err = CopyFileTool.Execute(ctx, `{"source":"`+srcFile+`","destination":"`+dstFile+`","confirm":true,"overwrite":false}`)
	assert.Error(t, err)

	// Overwrite with overwrite=true
	_, err = CopyFileTool.Execute(ctx, `{"source":"`+srcFile+`","destination":"`+dstFile+`","confirm":true,"overwrite":true}`)
	assert.NoError(t, err)
}

package declarations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestAppendFileTool_Metadata(t *testing.T) {
	assert.Equal(t, "append_file_tool", AppendFileTool.GetName())
	assert.Contains(t, AppendFileTool.GetCategories(), tool.FilesystemCategory)
}

func TestAppendFileTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, AppendFileTool, approvals)

	// confirm = false
	_, err := AppendFileTool.Execute(ctx, `{"path":"/tmp/test.txt","content":"hello","confirm":false}`)
	assert.Error(t, err)
	assert.Len(t, approvals.calls, 0)

	// empty path
	_, err = AppendFileTool.Execute(ctx, `{"path":"","content":"hello","confirm":true}`)
	assert.Error(t, err)
	assert.Len(t, approvals.calls, 0)
}

func TestAppendFileTool_DeniedByUser(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: false}
	withBindings(t, AppendFileTool, approvals)

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "denied.txt")

	_, err := AppendFileTool.Execute(ctx, `{"path":"`+filePath+`","content":"hello","confirm":true}`)
	assert.Error(t, err)
	assert.Len(t, approvals.calls, 1)
}

func TestAppendFileTool_Success(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, AppendFileTool, approvals)

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "append.txt")

	// 1. Create file by appending first chunk
	rawOut, err := AppendFileTool.Execute(ctx, `{"path":"`+filePath+`","content":"First line\n","confirm":true}`)
	assert.NoError(t, err)
	var out1 AppendFileOutput
	err = json.Unmarshal([]byte(rawOut), &out1)
	assert.NoError(t, err)
	assert.Equal(t, len("First line\n"), out1.Bytes)

	// 2. Append second chunk
	rawOut2, err := AppendFileTool.Execute(ctx, `{"path":"`+filePath+`","content":"Second line\n","confirm":true}`)
	assert.NoError(t, err)
	var out2 AppendFileOutput
	err = json.Unmarshal([]byte(rawOut2), &out2)
	assert.NoError(t, err)
	assert.Equal(t, len("Second line\n"), out2.Bytes)

	// Verify content on disk
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "First line\nSecond line\n", string(content))
}

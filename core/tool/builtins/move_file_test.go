package declarations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestMoveFileTool_Metadata(t *testing.T) {
	assert.Equal(t, "move_file_tool", MoveFileTool.GetName())
	assert.Contains(t, MoveFileTool.GetCategories(), tool.FilesystemCategory)
}

func TestMoveFileTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, MoveFileTool, approvals)

	_, err := MoveFileTool.Execute(ctx, `{"source":"a","destination":"b","confirm":false}`)
	assert.Error(t, err)

	_, err = MoveFileTool.Execute(ctx, `{"source":"","destination":"b","confirm":true}`)
	assert.Error(t, err)
}

func TestMoveFileTool_Success(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, MoveFileTool, approvals)

	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "orig.txt")
	dstFile := filepath.Join(tmpDir, "renamed.txt")

	err := os.WriteFile(srcFile, []byte("data to move"), 0644)
	assert.NoError(t, err)

	rawOut, err := MoveFileTool.Execute(ctx, `{"source":"`+srcFile+`","destination":"`+dstFile+`","confirm":true}`)
	assert.NoError(t, err)

	var out MoveFileOutput
	err = json.Unmarshal([]byte(rawOut), &out)
	assert.NoError(t, err)

	// Check that source no longer exists and destination exists with content
	_, err = os.Stat(srcFile)
	assert.True(t, os.IsNotExist(err))

	content, err := os.ReadFile(dstFile)
	assert.NoError(t, err)
	assert.Equal(t, "data to move", string(content))
}

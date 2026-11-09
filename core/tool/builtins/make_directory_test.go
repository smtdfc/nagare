package declarations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestMakeDirectoryTool_Metadata(t *testing.T) {
	assert.Equal(t, "make_directory_tool", MakeDirectoryTool.GetName())
	assert.Contains(t, MakeDirectoryTool.GetCategories(), tool.FilesystemCategory)
}

func TestMakeDirectoryTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, MakeDirectoryTool, approvals)

	_, err := MakeDirectoryTool.Execute(ctx, `{"path":"/tmp/a","confirm":false}`)
	assert.Error(t, err)

	_, err = MakeDirectoryTool.Execute(ctx, `{"path":"","confirm":true}`)
	assert.Error(t, err)
}

func TestMakeDirectoryTool_Success(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, MakeDirectoryTool, approvals)

	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "level1", "level2")

	rawOut, err := MakeDirectoryTool.Execute(ctx, `{"path":"`+nestedDir+`","confirm":true}`)
	assert.NoError(t, err)

	var out MakeDirectoryOutput
	err = json.Unmarshal([]byte(rawOut), &out)
	assert.NoError(t, err)

	info, err := os.Stat(nestedDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

package declarations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestFindFilesTool_Metadata(t *testing.T) {
	assert.Equal(t, "find_files_tool", FindFilesTool.GetName())
	assert.Contains(t, FindFilesTool.GetCategories(), tool.FilesystemCategory)
}

func TestFindFilesTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	_, err := FindFilesTool.Execute(ctx, `{"pattern":""}`)
	assert.Error(t, err)

	_, err = FindFilesTool.Execute(ctx, `{"directory":"/non/existent/dir/xyz", "pattern":"*.txt"}`)
	assert.Error(t, err)
}

func TestFindFilesTool_Success(t *testing.T) {
	ctx := approvalTestContext()
	tmpDir := t.TempDir()

	// Create test structure:
	// tmpDir/
	//   a.txt
	//   b.log
	//   sub/
	//     c.txt
	//     deep/
	//       d.txt
	assert.NoError(t, os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("a"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(tmpDir, "b.log"), []byte("b"), 0644))
	subDir := filepath.Join(tmpDir, "sub")
	assert.NoError(t, os.MkdirAll(filepath.Join(subDir, "deep"), 0755))
	assert.NoError(t, os.WriteFile(filepath.Join(subDir, "c.txt"), []byte("c"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(subDir, "deep", "d.txt"), []byte("d"), 0644))

	// Find all *.txt
	rawOut, err := FindFilesTool.Execute(ctx, `{"directory":"`+tmpDir+`","pattern":"*.txt"}`)
	assert.NoError(t, err)

	var out FindFilesOutput
	err = json.Unmarshal([]byte(rawOut), &out)
	assert.NoError(t, err)
	assert.Equal(t, 3, out.Total)
	assert.Len(t, out.Matches, 3)

	// Find all *.log
	rawOutLog, err := FindFilesTool.Execute(ctx, `{"directory":"`+tmpDir+`","pattern":"*.log"}`)
	assert.NoError(t, err)
	var outLog FindFilesOutput
	err = json.Unmarshal([]byte(rawOutLog), &outLog)
	assert.NoError(t, err)
	assert.Equal(t, 1, outLog.Total)
}

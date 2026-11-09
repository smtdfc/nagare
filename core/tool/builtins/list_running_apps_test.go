package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestListRunningAppsTool_Metadata(t *testing.T) {
	assert.Equal(t, "list_running_apps_tool", ListRunningAppsTool.GetName())
	assert.NotEmpty(t, ListRunningAppsTool.GetDescription())
	assert.Contains(t, ListRunningAppsTool.GetCategories(), tool.ComputerUse)
}

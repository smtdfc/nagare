package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestGetUITreeTool_Metadata(t *testing.T) {
	assert.Equal(t, "get_ui_tree_tool", GetUITreeTool.GetName())
	assert.NotEmpty(t, GetUITreeTool.GetDescription())
	assert.Contains(t, GetUITreeTool.GetCategories(), tool.ComputerUse)
}

func TestGetUITreeTool_Validation(t *testing.T) {
	execCtx := &core_context.ExecuteContext{
		Context: context.Background(),
	}

	_, err := GetUITreeTool.Execute(execCtx, `{"instance_id": ""}`)
	assert.Error(t, err)

	_, err = GetUITreeTool.Execute(execCtx, `{}`)
	assert.Error(t, err)
}

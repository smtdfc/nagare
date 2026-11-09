package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestOpenAppTool_Metadata(t *testing.T) {
	assert.Equal(t, "open_app_tool", OpenAppTool.GetName())
	assert.NotEmpty(t, OpenAppTool.GetDescription())
	assert.Contains(t, OpenAppTool.GetCategories(), tool.ComputerUse)
}

func TestOpenAppTool_Validation(t *testing.T) {
	execCtx := &core_context.ExecuteContext{
		Context: context.Background(),
	}

	_, err := OpenAppTool.Execute(execCtx, `{"name": ""}`)
	assert.Error(t, err)

	_, err = OpenAppTool.Execute(execCtx, `{}`)
	assert.Error(t, err)
}

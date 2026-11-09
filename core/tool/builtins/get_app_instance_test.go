package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestGetAppInstanceTool_Metadata(t *testing.T) {
	assert.Equal(t, "get_app_instance_tool", GetAppInstanceTool.GetName())
	assert.NotEmpty(t, GetAppInstanceTool.GetDescription())
	assert.Contains(t, GetAppInstanceTool.GetCategories(), tool.ComputerUse)
}

func TestGetAppInstanceTool_Validation(t *testing.T) {
	execCtx := &core_context.ExecuteContext{
		Context: context.Background(),
	}

	_, err := GetAppInstanceTool.Execute(execCtx, `{"name": ""}`)
	assert.Error(t, err)

	_, err = GetAppInstanceTool.Execute(execCtx, `{}`)
	assert.Error(t, err)
}

package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestWindowControlTool_Metadata(t *testing.T) {
	assert.Equal(t, "window_control_tool", WindowControlTool.GetName())
	assert.NotEmpty(t, WindowControlTool.GetDescription())
	assert.Contains(t, WindowControlTool.GetCategories(), tool.ComputerUse)
}

func TestWindowControlTool_Validation(t *testing.T) {
	execCtx := &core_context.ExecuteContext{
		Context: context.Background(),
	}

	_, err := WindowControlTool.Execute(execCtx, `{"instance_id": ""}`)
	assert.Error(t, err)

	_, err = WindowControlTool.Execute(execCtx, `{"instance_id": "app-1", "action": ""}`)
	assert.Error(t, err)

	_, err = WindowControlTool.Execute(execCtx, `{}`)
	assert.Error(t, err)
}

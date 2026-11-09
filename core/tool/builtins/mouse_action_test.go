package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestMouseActionTool_Metadata(t *testing.T) {
	assert.Equal(t, "mouse_action_tool", MouseActionTool.GetName())
	assert.NotEmpty(t, MouseActionTool.GetDescription())
	assert.Contains(t, MouseActionTool.GetCategories(), tool.ComputerUse)
}

func TestMouseActionTool_Validation(t *testing.T) {
	execCtx := &core_context.ExecuteContext{
		Context: context.Background(),
	}

	_, err := MouseActionTool.Execute(execCtx, `{"action": ""}`)
	assert.Error(t, err)

	_, err = MouseActionTool.Execute(execCtx, `{}`)
	assert.Error(t, err)
}

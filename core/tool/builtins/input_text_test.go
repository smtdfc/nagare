package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestInputTextTool_Metadata(t *testing.T) {
	assert.Equal(t, "input_text_tool", InputTextTool.GetName())
	assert.NotEmpty(t, InputTextTool.GetDescription())
	assert.Contains(t, InputTextTool.GetCategories(), tool.ComputerUse)
}

func TestInputTextTool_Validation(t *testing.T) {
	execCtx := &core_context.ExecuteContext{
		Context: context.Background(),
	}

	_, err := InputTextTool.Execute(execCtx, `{"text": ""}`)
	assert.Error(t, err)

	_, err = InputTextTool.Execute(execCtx, `{}`)
	assert.Error(t, err)
}

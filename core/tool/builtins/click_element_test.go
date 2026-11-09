package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestClickElementTool_Metadata(t *testing.T) {
	assert.Equal(t, "click_element_tool", ClickElementTool.GetName())
	assert.NotEmpty(t, ClickElementTool.GetDescription())
	assert.Contains(t, ClickElementTool.GetCategories(), tool.ComputerUse)
}

func TestClickElementTool_Validation(t *testing.T) {
	execCtx := &core_context.ExecuteContext{
		Context: context.Background(),
	}

	_, err := ClickElementTool.Execute(execCtx, `{"instance_id": ""}`)
	assert.Error(t, err)

	_, err = ClickElementTool.Execute(execCtx, `{"instance_id": "app-1", "ref_id": ""}`)
	assert.Error(t, err)

	_, err = ClickElementTool.Execute(execCtx, `{}`)
	assert.Error(t, err)
}

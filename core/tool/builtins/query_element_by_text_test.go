package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestQueryElementByTextTool_Metadata(t *testing.T) {
	assert.Equal(t, "query_element_by_text_tool", QueryElementByTextTool.GetName())
	assert.NotEmpty(t, QueryElementByTextTool.GetDescription())
	assert.Contains(t, QueryElementByTextTool.GetCategories(), tool.ComputerUse)
}

func TestQueryElementByTextTool_Validation(t *testing.T) {
	execCtx := &core_context.ExecuteContext{
		Context: context.Background(),
	}

	_, err := QueryElementByTextTool.Execute(execCtx, `{"instance_id": "", "text": ""}`)
	assert.Error(t, err)

	_, err = QueryElementByTextTool.Execute(execCtx, `{"instance_id": "app-1", "text": ""}`)
	assert.Error(t, err)

	_, err = QueryElementByTextTool.Execute(execCtx, `{"instance_id": "", "text": "OK"}`)
	assert.Error(t, err)

	_, err = QueryElementByTextTool.Execute(execCtx, `{}`)
	assert.Error(t, err)
}

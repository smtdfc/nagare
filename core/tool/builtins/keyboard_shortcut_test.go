package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestKeyboardShortcutTool_Metadata(t *testing.T) {
	assert.Equal(t, "keyboard_shortcut_tool", KeyboardShortcutTool.GetName())
	assert.NotEmpty(t, KeyboardShortcutTool.GetDescription())
	assert.Contains(t, KeyboardShortcutTool.GetCategories(), tool.ComputerUse)
}

func TestKeyboardShortcutTool_Validation(t *testing.T) {
	execCtx := &core_context.ExecuteContext{
		Context: context.Background(),
	}

	_, err := KeyboardShortcutTool.Execute(execCtx, `{"key": ""}`)
	assert.Error(t, err)

	_, err = KeyboardShortcutTool.Execute(execCtx, `{}`)
	assert.Error(t, err)
}

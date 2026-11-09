package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type KeyboardShortcutInput struct {
	Key       string   `json:"key" jsonschema_description:"Main key to press, such as 'Enter', 'Escape', 'Tab', 'Space', 'Backspace', 'Delete', 'ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'Home', 'End', 'PageUp', 'PageDown', 'F1'-'F12', or any character 'a'-'z', '0'-'9'."`
	Modifiers []string `json:"modifiers,omitempty" jsonschema_description:"Optional list of modifier keys to hold down during key press. Supported values: 'ctrl', 'alt', 'shift', 'meta' (Command/Windows key). E.g. ['ctrl'] for Ctrl+C, ['alt'] for Alt+Tab."`
}

type KeyboardShortcutOutput struct {
	Success bool `json:"success" jsonschema_description:"Indicates whether the keyboard key/shortcut was successfully simulated."`
}

var KeyboardShortcutTool = tool.DefineTool(
	"keyboard_shortcut_tool",
	"Simulate native OS keyboard key presses or hotkey combinations (chords) such as Ctrl+C, Ctrl+V, Alt+Tab, Escape, Enter, or Function keys.",
	func(ctx *context.ExecuteContext, args *KeyboardShortcutInput, bindings tool.Bindings) (*KeyboardShortcutOutput, error) {
		if args == nil || args.Key == "" {
			return nil, fmt.Errorf("key cannot be empty")
		}

		err := bindings.GetComputerUseOrchestrator().KeyboardShortcut(ctx, args.Key, args.Modifiers)
		if err != nil {
			return nil, err
		}

		return &KeyboardShortcutOutput{Success: true}, nil
	},
	[]string{tool.ComputerUse},
)

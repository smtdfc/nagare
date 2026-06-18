package chat

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// View implements [tui.Page].
func (m *ChatPage) View() tea.View {
	return tea.NewView(
		fmt.Sprintf(
			"%s\n\n%s",
			m.viewport.View(),
			m.textarea.View(),
		),
	)
}

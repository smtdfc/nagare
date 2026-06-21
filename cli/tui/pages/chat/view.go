package chat

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// View implements [tui.Page].
func (m *ChatPage) View() tea.View {
	view := tea.NewView(
		fmt.Sprintf(
			"%s\n\n%s",
			m.viewport.View(),
			m.textarea.View(),
		),
	)

	view.AltScreen = true

	return view
}

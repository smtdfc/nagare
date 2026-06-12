package chat

import (
	tea "charm.land/bubbletea/v2"
)

// View implements [tui.Page].
func (m *ChatPage) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true
	return view
}

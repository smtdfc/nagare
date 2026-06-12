package provider_setting_page

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// View implements [tui.Page].
func (m *Page) View() tea.View {
	view := tea.NewView(lipgloss.NewStyle().Margin(2).Render(m.list.View()))
	view.AltScreen = true
	return view
}

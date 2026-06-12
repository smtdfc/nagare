package chat

import (
	tea "charm.land/bubbletea/v2"
)

func (m *ChatPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		//

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
			// case "enter":
			// 	selectedItem := m.list.SelectedItem()
			// 	if item, ok := selectedItem.(components.ListItem); ok {
			// 		return m, func() tea.Msg {
			// 			return router.ChangePageMsg{
			// 				Target: item.FilterValue(),
			// 			}
			// 		}
			// 	}
		}

	}

	return m, cmd
}

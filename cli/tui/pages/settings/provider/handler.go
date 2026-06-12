package provider_setting_page

import (
	tea "charm.land/bubbletea/v2"

	"github.com/smtdfc/nagare/cli/tui/components"
	"github.com/smtdfc/nagare/cli/tui/router"
)

func (m *Page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-1)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selectedItem := m.list.SelectedItem()
			if item, ok := selectedItem.(components.ListItem[string]); ok {
				return m, func() tea.Msg {
					return router.ChangePageMsg{
						Target:  item.Value(),
						Refresh: true,
					}
				}
			}
		}

	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

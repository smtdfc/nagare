package plugin_setting_page

import (
	"sort"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/smtdfc/nagare/cli/tui/router"
	"github.com/smtdfc/nagare/core/config"
)

type ListPage struct {
	tea.Model
	table  table.Model
	config *config.Config
	width  int
	height int
}

func NewListPage(conf *config.Config) *ListPage {
	m := &ListPage{config: conf}
	m.Refresh()
	return m
}

func (m *ListPage) GetName() string {
	return "settings:plugin:list"
}

func (m *ListPage) Init() tea.Cmd {
	return nil
}

func (m *ListPage) Refresh() {
	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Name", Width: 40},
		{Title: "Author", Width: 40},
		{Title: "Version", Width: 10},
	}

	var names []string
	for name := range m.config.Plugins {
		names = append(names, name)
	}
	sort.Strings(names)

	var rows []table.Row
	for _, name := range names {
		p := m.config.Plugins[name]

		rows = append(rows, table.Row{
			p.Id,
			p.Name,
			p.Author,
			p.Version,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithWidth(100),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(s)

	if m.width > 0 {
		t.SetWidth(m.width)
	}
	if m.height > 0 {
		t.SetHeight(m.height - 8)
	} else {
		t.SetHeight(10)
	}

	m.table = t
}

func (m *ListPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height - 8)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg {
				return router.ChangePageMsg{
					Target:  "settings:plugin",
					Refresh: true,
				}
			}
		case "ctrl+c":
			return m, tea.Quit

		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *ListPage) View() tea.View {
	doc := lipgloss.NewStyle().Margin(2)

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Active Plugin")
	helpMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press [Esc/q] to go back")

	viewContent := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		m.table.View(),
		"",
		helpMsg,
	)

	view := tea.NewView(doc.Render(viewContent))
	view.AltScreen = true
	return view
}

package provider_setting_page

import (
	"sort"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/smtdfc/nagare/cli/tui/router"
	"github.com/smtdfc/nagare/core/config"
)

type RemovePage struct {
	tea.Model
	table  table.Model
	config *config.Config
	width  int
	height int
}

func NewRemovePage(conf *config.Config) *RemovePage {
	m := &RemovePage{config: conf}
	m.Refresh()
	return m
}

func (m *RemovePage) GetName() string {
	return "settings:provider:remove"
}

func (m *RemovePage) Init() tea.Cmd {
	return nil
}

func (m *RemovePage) Refresh() {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Base URL", Width: 40},
		{Title: "Model Name", Width: 20},
		{Title: "Compatible", Width: 15},
	}

	var names []string
	for name := range m.config.Providers {
		names = append(names, name)
	}
	sort.Strings(names)

	var rows []table.Row
	for _, name := range names {
		p := m.config.Providers[name]
		rows = append(rows, table.Row{name, p.BaseURL, p.ModelName, string(p.Compatible)})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithWidth(100),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("160")).Bold(true)
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

func (m *RemovePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				return router.ChangePageMsg{Target: "settings:provider", Refresh: true}
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected := m.table.SelectedRow()
			if len(selected) > 0 {
				val := selected[0]
				delete(m.config.Providers, val)
				if m.config.CurrentProvider == val {
					m.config.CurrentProvider = ""
					for remainingName := range m.config.Providers {
						m.config.CurrentProvider = remainingName
						break
					}
				}
				_ = config.SaveConfig(m.config)
				return m, func() tea.Msg {
					return router.ChangePageMsg{Target: "settings:provider", Refresh: true}
				}
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *RemovePage) View() tea.View {
	doc := lipgloss.NewStyle().Margin(2)
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("160")).Render("Select Provider to Remove")
	helpMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press [Enter] to delete • Press [Esc/q] to go back")

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

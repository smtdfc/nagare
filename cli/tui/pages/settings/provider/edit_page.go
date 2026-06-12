package provider_setting_page

import (
	"fmt"
	"sort"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/smtdfc/nagare/cli/tui/router"
	"github.com/smtdfc/nagare/core/config"
)

type EditPage struct {
	tea.Model
	config       *config.Config
	table        table.Model
	inputs       []textinput.Model
	focusIndex   int
	selectedName string
	width        int
	height       int
}

func NewEditPage(conf *config.Config) *EditPage {
	m := &EditPage{config: conf}
	m.Refresh()
	return m
}

func (m *EditPage) GetName() string {
	return "settings:provider:edit"
}

func (m *EditPage) Init() tea.Cmd {
	return textinput.Blink
}

func (m *EditPage) Refresh() {
	m.selectedName = ""
	m.focusIndex = 0

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
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(true)
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

func (m *EditPage) startEditing(name string) {
	m.selectedName = name
	m.focusIndex = 0
	p := m.config.Providers[name]

	m.inputs = make([]textinput.Model, 5)

	var t textinput.Model

	t = textinput.New()
	t.Prompt = "Name:        "
	t.SetValue(p.Name)
	t.Focus()
	t.SetWidth(40)
	m.inputs[0] = t

	t = textinput.New()
	t.Prompt = "Base URL:    "
	t.SetValue(p.BaseURL)
	t.SetWidth(40)
	m.inputs[1] = t

	t = textinput.New()
	t.Prompt = "API Key:     "
	t.SetValue(p.APIKey)
	t.EchoMode = textinput.EchoPassword
	t.EchoCharacter = '•'
	t.SetWidth(40)
	m.inputs[2] = t

	t = textinput.New()
	t.Prompt = "Model Name:  "
	t.SetValue(p.ModelName)
	t.SetWidth(40)
	m.inputs[3] = t

	t = textinput.New()
	t.Prompt = "Compatible:  "
	t.SetValue(string(p.Compatible))
	t.SetWidth(40)
	m.inputs[4] = t
}

func (m *EditPage) saveProvider() {
	name := m.inputs[0].Value()
	baseURL := m.inputs[1].Value()
	apiKey := m.inputs[2].Value()
	modelName := m.inputs[3].Value()
	comp := config.ProviderCompatible(m.inputs[4].Value())

	if name == "" {
		return
	}

	if m.selectedName != name {
		delete(m.config.Providers, m.selectedName)
		if m.config.CurrentProvider == m.selectedName {
			m.config.CurrentProvider = name
		}
	}

	p := config.Provider{
		Name:       name,
		BaseURL:    baseURL,
		APIKey:     apiKey,
		ModelName:  modelName,
		Compatible: comp,
		Enabled:    true,
	}

	m.config.Providers[name] = p
	_ = config.SaveConfig(m.config)
}

func (m *EditPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height - 8)
	}

	if m.selectedName == "" {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q":
				return m, func() tea.Msg { return router.ChangePageMsg{Target: "settings:provider", Refresh: true} }
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				selected := m.table.SelectedRow()
				if len(selected) > 0 {
					m.startEditing(selected[0])
					return m, nil
				}
			}
		}
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.selectedName = ""
			m.Refresh()
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "tab", "shift+tab", "up", "down":
			s := msg.String()
			backward := s == "shift+tab" || s == "up"

			if m.focusIndex < len(m.inputs) {
				m.inputs[m.focusIndex].Blur()
			}

			if backward {
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs) + 1
				}
			} else {
				m.focusIndex++
				if m.focusIndex > len(m.inputs)+1 {
					m.focusIndex = 0
				}
			}

			if m.focusIndex < len(m.inputs) {
				m.inputs[m.focusIndex].Focus()
			}
			return m, nil

		case "enter":
			if m.focusIndex == len(m.inputs) {
				m.saveProvider()
				return m, func() tea.Msg {
					return router.ChangePageMsg{Target: "settings:provider", Refresh: true}
				}
			} else if m.focusIndex == len(m.inputs)+1 {
				m.selectedName = ""
				m.Refresh()
				return m, nil
			} else {
				if m.focusIndex < len(m.inputs) {
					m.inputs[m.focusIndex].Blur()
					m.focusIndex++
					if m.focusIndex < len(m.inputs) {
						m.inputs[m.focusIndex].Focus()
					}
				}
				return m, nil
			}
		}
	}

	if m.focusIndex < len(m.inputs) {
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	}

	return m, cmd
}

func (m *EditPage) View() tea.View {
	doc := lipgloss.NewStyle().Margin(2)

	if m.selectedName == "" {
		helpMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press [Enter] to edit • [Esc/q] to go back")
		content := lipgloss.JoinVertical(lipgloss.Left, m.table.View(), "", helpMsg)
		view := tea.NewView(doc.Render(content))
		view.AltScreen = true
		return view
	}

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render(fmt.Sprintf("Edit Provider: %s", m.selectedName))

	var views []string
	for i := range m.inputs {
		views = append(views, m.inputs[i].View())
	}

	submitBtn := " [ Save ] "
	cancelBtn := " [ Cancel ] "

	if m.focusIndex == len(m.inputs) {
		submitBtn = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("205")).Foreground(lipgloss.Color("0")).Render(" [ Save ] ")
	} else if m.focusIndex == len(m.inputs)+1 {
		cancelBtn = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("240")).Foreground(lipgloss.Color("0")).Render(" [ Cancel ] ")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Left, submitBtn, "  ", cancelBtn)

	formView := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		views[0], views[1], views[2], views[3], views[4],
		"",
		buttons,
	)

	view := tea.NewView(doc.Render(formView))
	view.AltScreen = true
	return view
}

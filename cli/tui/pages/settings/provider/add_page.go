package provider_setting_page

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/smtdfc/nagare/cli/tui/router"
	"github.com/smtdfc/nagare/core/config"
)

type AddPage struct {
	tea.Model
	config     *config.Config
	inputs     []textinput.Model
	focusIndex int
}

func NewAddPage(conf *config.Config) *AddPage {
	m := &AddPage{config: conf}
	m.Refresh()
	return m
}

func (m *AddPage) GetName() string {
	return "settings:provider:add"
}

func (m *AddPage) Init() tea.Cmd {
	return textinput.Blink
}

func (m *AddPage) Refresh() {
	m.inputs = make([]textinput.Model, 5)
	m.focusIndex = 0

	var t textinput.Model

	t = textinput.New()
	t.Placeholder = "e.g., Groq"
	t.Focus()
	t.Prompt = "Name:        "
	t.SetWidth(40)
	m.inputs[0] = t

	t = textinput.New()
	t.Placeholder = "e.g., https://api.groq.com/openai/v1"
	t.Prompt = "Base URL:    "
	t.SetWidth(40)
	m.inputs[1] = t

	t = textinput.New()
	t.Placeholder = "e.g., gsk_..."
	t.EchoMode = textinput.EchoPassword
	t.EchoCharacter = '•'
	t.Prompt = "API Key:     "
	t.SetWidth(40)
	m.inputs[2] = t

	t = textinput.New()
	t.Placeholder = "e.g., llama3-8b-8192"
	t.Prompt = "Model Name:  "
	t.SetWidth(40)
	m.inputs[3] = t

	t = textinput.New()
	t.Placeholder = "openai"
	t.SetValue("openai")
	t.Prompt = "Compatible:  "
	t.SetWidth(40)
	m.inputs[4] = t
}

func (m *AddPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
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
					return router.ChangePageMsg{
						Target:  "settings:provider",
						Refresh: true,
					}
				}
			} else if m.focusIndex == len(m.inputs)+1 {
				return m, func() tea.Msg {
					return router.ChangePageMsg{
						Target:  "settings:provider",
						Refresh: true,
					}
				}
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

func (m *AddPage) saveProvider() {
	name := m.inputs[0].Value()
	baseURL := m.inputs[1].Value()
	apiKey := m.inputs[2].Value()
	modelName := m.inputs[3].Value()
	comp := config.ProviderCompatible(m.inputs[4].Value())

	if name == "" {
		return
	}

	p := config.Provider{
		Name:       name,
		BaseURL:    baseURL,
		APIKey:     apiKey,
		ModelName:  modelName,
		Compatible: comp,
		Enabled:    true,
	}

	if m.config.Providers == nil {
		m.config.Providers = make(map[string]config.Provider)
	}
	m.config.Providers[name] = p

	if m.config.CurrentProvider == "" {
		m.config.CurrentProvider = name
	}

	_ = config.SaveConfig(m.config)
}

func (m *AddPage) View() tea.View {
	doc := lipgloss.NewStyle().Margin(2)
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Add New LLM Provider")

	var views []string
	for i := range m.inputs {
		views = append(views, m.inputs[i].View())
	}

	submitBtn := " [ Submit ] "
	cancelBtn := " [ Cancel ] "

	if m.focusIndex == len(m.inputs) {
		submitBtn = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("205")).Foreground(lipgloss.Color("0")).Render(" [ Submit ] ")
	} else if m.focusIndex == len(m.inputs)+1 {
		cancelBtn = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("240")).Foreground(lipgloss.Color("0")).Render(" [ Cancel ] ")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Left, submitBtn, "  ", cancelBtn)

	formView := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		views[0],
		views[1],
		views[2],
		views[3],
		views[4],
		"",
		buttons,
	)

	view := tea.NewView(doc.Render(formView))
	view.AltScreen = true
	return view
}

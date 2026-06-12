package main_ui

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/smtdfc/nagare/cli/tui/components"
)

type MainPage struct {
	tea.Model
	list list.Model
}

// Refresh implements [router.Page].
func (m *MainPage) Refresh() {
	panic("unimplemented")
}

// GetName implements [tui.Page].
func (m *MainPage) GetName() string {
	return "main"
}

// Init implements [tui.Page].
func (m *MainPage) Init() tea.Cmd {
	return nil
}

func NewMainPage() *MainPage {
	items := []list.Item{
		components.NewListItem("Chat", "Chat with Nagare", "chat"),
		components.NewListItem("Settings", "Settings for Nagare", "settings"),
		components.NewListItem("Plugin", "Manage your plugins", "settings:plugin"),
	}

	l := list.New(items, list.NewDefaultDelegate(), 100, 20)
	l.Title = "Nagare Agent"
	return &MainPage{list: l}
}

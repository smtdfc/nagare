package chat

import (
	tea "charm.land/bubbletea/v2"
)

type ChatPage struct {
	tea.Model
}

// GetName implements [tui.Page].
func (m *ChatPage) GetName() string {
	return "chat"
}

// Init implements [tui.Page].
func (m *ChatPage) Init() tea.Cmd {
	return nil
}

func NewMainPage() *ChatPage {
	return &ChatPage{}
}

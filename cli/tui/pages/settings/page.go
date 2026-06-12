package settings

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/smtdfc/nagare/cli/tui/components"
	"github.com/smtdfc/nagare/core/config"
)

type Page struct {
	tea.Model
	list   list.Model
	config *config.Config
	width  int
	height int
}

func BuildList(conf *config.Config) list.Model {
	items := []list.Item{
		components.NewListItem("Mode", fmt.Sprintf("Select mode for Nagare (Current: %s)", conf.CurrentMode.String()), "settings:mode"),
		components.NewListItem("Provider", fmt.Sprintf("Manage your LLM providers (Added: %d provider)", len(conf.Providers)), "settings:provider"),
		components.NewListItem("Plugin", fmt.Sprintf("Manage your Plugin (Added: %d plugin)", len(conf.Plugins)), "settings:plugin"),
		components.NewListItem("Back", "", "main"),
	}

	l := list.New(items, list.NewDefaultDelegate(), 100, 20)
	l.Title = "Nagare Agent Settings"

	return l
}

// Refresh implements [router.Page].
func (m *Page) Refresh() {
	m.list = BuildList(m.config)
	if m.width > 0 && m.height > 0 {
		m.list.SetSize(m.width, m.height-1)
	}
}

// GetName implements [tui.Page].
func (m *Page) GetName() string {
	return "settings"
}

// Init implements [tui.Page].
func (m *Page) Init() tea.Cmd {
	return nil
}

func NewPage(conf *config.Config) *Page {
	return &Page{list: BuildList(conf), config: conf}
}

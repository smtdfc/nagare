package provider_setting_page

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/smtdfc/nagare/cli/tui/components"
	"github.com/smtdfc/nagare/core/config"
)

type Page struct {
	tea.Model
	list   list.Model
	config *config.Config
}

// Refresh implements [router.Page].
func (m *Page) Refresh() {
	// Reset selection to the first item on entry
	m.list.Select(0)
}

// GetName implements [tui.Page].
func (m *Page) GetName() string {
	return "settings:provider"
}

// Init implements [tui.Page].
func (m *Page) Init() tea.Cmd {
	return nil
}

func NewPage(conf *config.Config) *Page {
	items := []list.Item{
		components.NewListItem("List", "Show list provider", "settings:provider:list"),
		components.NewListItem("Add", "Add provider", "settings:provider:add"),
		components.NewListItem("Edit", "Edit provider", "settings:provider:edit"),
		components.NewListItem("Remove", "Remove provider", "settings:provider:remove"),
		components.NewListItem("Back", "Return to main settings", "settings"),
	}

	l := list.New(items, list.NewDefaultDelegate(), 100, 20)
	l.Title = "Manage Providers"
	return &Page{list: l, config: conf}
}

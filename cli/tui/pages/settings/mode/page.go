package mode_setting_page

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
	for i, item := range m.list.Items() {
		if li, ok := item.(components.ListItem[config.Mode]); ok {
			if li.Value() == m.config.CurrentMode {
				m.list.Select(i)
				break
			}
		}
	}
}

// GetName implements [tui.Page].
func (m *Page) GetName() string {
	return "settings:mode"
}

// Init implements [tui.Page].
func (m *Page) Init() tea.Cmd {
	return nil
}

func NewPage(conf *config.Config) *Page {
	items := []list.Item{
		components.NewListItem("Provider", "Nagare will use a model direct from a LLM provider", config.PROVIDER_MODE),
		components.NewListItem("Proxy", "Nagare will use a model from an AI Proxy/Router", config.PROXY_MODE),
	}

	l := list.New(items, list.NewDefaultDelegate(), 100, 20)
	l.Title = "Select a mode: "
	return &Page{list: l, config: conf}
}

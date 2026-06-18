package tui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/smtdfc/nagare/cli/tui/pages/chat"
	"github.com/smtdfc/nagare/cli/tui/pages/main_ui"
	"github.com/smtdfc/nagare/cli/tui/pages/settings"
	mode_setting_page "github.com/smtdfc/nagare/cli/tui/pages/settings/mode"
	plugin_setting_page "github.com/smtdfc/nagare/cli/tui/pages/settings/plugin"
	provider_setting_page "github.com/smtdfc/nagare/cli/tui/pages/settings/provider"
	"github.com/smtdfc/nagare/cli/tui/router"
	"github.com/smtdfc/nagare/core"
	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/config"
)

type RootTUIModel struct {
	Router     *router.TUIRouter
	Config     *config.Config
	SessionMgr *agent.SessionManager
	AgentPool  *agent.AgentPool
}

// Init implements [tea.Model].
func (r *RootTUIModel) Init() tea.Cmd {
	return nil
}

// Update implements [tea.Model].
func (r *RootTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := r.Router.Pages[r.Router.CurrentPage]
	newModel, cmd := current.Update(msg)
	r.Router.Pages[r.Router.CurrentPage] = newModel.(router.Page)

	switch msg := msg.(type) {
	case router.ChangePageMsg:
		if r.Router.HasPage(msg.Target) {
			r.Router.CurrentPage = msg.Target
			if msg.Refresh {
				r.Router.Pages[r.Router.CurrentPage].Refresh()
			}
		} else {
			fmt.Printf("Page %s not found. \n", msg.Target)
		}
		return r, nil
	}

	return r, cmd
}

// View implements [tea.Model].
func (r *RootTUIModel) View() tea.View {
	current := r.Router.Pages[r.Router.CurrentPage]
	return current.View()
}

func NewRootTUI(conf *config.Config) {

	agentPool, sessionMgr := core.InitAgent(conf)
	model := &RootTUIModel{
		Router:     router.NewTUIRouter(),
		Config:     conf,
		SessionMgr: sessionMgr,
		AgentPool:  agentPool,
	}
	p := tea.NewProgram(model)

	model.Router.Pages["main"] = main_ui.NewMainPage()
	model.Router.Pages["chat"] = chat.NewPage(sessionMgr, agentPool)
	model.Router.Pages["settings"] = settings.NewPage(conf)

	// Mode settings
	model.Router.Pages["settings:mode"] = mode_setting_page.NewPage(conf)

	// Provider settings
	model.Router.Pages["settings:provider"] = provider_setting_page.NewPage(conf)
	model.Router.Pages["settings:provider:list"] = provider_setting_page.NewListPage(conf)
	model.Router.Pages["settings:provider:add"] = provider_setting_page.NewAddPage(conf)
	model.Router.Pages["settings:provider:edit"] = provider_setting_page.NewEditPage(conf)
	model.Router.Pages["settings:provider:remove"] = provider_setting_page.NewRemovePage(conf)

	// Plugin settings
	model.Router.Pages["settings:plugin"] = plugin_setting_page.NewPage(conf)
	model.Router.Pages["settings:plugin:list"] = plugin_setting_page.NewListPage(conf)
	model.Router.Pages["settings:plugin:remove"] = plugin_setting_page.NewRemovePage(conf)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

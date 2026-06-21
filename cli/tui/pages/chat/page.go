package chat

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/messages"
)

type ChatPage struct {
	tea.Model
	agentPool  *agent.AgentPool
	sessionMgr *agent.SessionManager
	sessionID  string

	viewport viewport.Model
	textarea textarea.Model

	currentMessageChunkType string
	textMessage             string
	reasoningMessage        string
	messages                []string
	err                     error

	currentStream <-chan messages.Message
	currentAgent  *agent.Agent
}

// GetName implements [tui.Page].
func (m *ChatPage) GetName() string {
	return "chat"
}

// Init implements [tui.Page].
func (m *ChatPage) Init() tea.Cmd {
	return textarea.Blink
}

func (m *ChatPage) Refresh() {}

func NewPage(sessionMgr *agent.SessionManager, agentPool *agent.AgentPool) *ChatPage {
	ta := textarea.New()
	ta.Placeholder = "Type a message... (Press Enter to send)"
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 500
	ta.SetWidth(100)
	ta.SetHeight(3)
	ta.ShowLineNumbers = false

	vp := viewport.New()
	vp.SetWidth(100)
	vp.SetHeight(20)
	vp.SetContent("Welcome to Nagare Agent Chat!")

	return &ChatPage{
		agentPool:               agentPool,
		sessionMgr:              sessionMgr,
		sessionID:               "1",
		textarea:                ta,
		viewport:                vp,
		messages:                []string{},
		currentMessageChunkType: "unknown",
	}
}

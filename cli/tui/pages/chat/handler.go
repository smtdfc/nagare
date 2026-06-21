package chat

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/smtdfc/nagare/cli/tui/router"
)

func (m *ChatPage) renderReasoningMessage(msg string) string {
	style := lipgloss.NewStyle().Italic(true).MarginLeft(1).Foreground(lipgloss.Color("#605959"))
	return style.Render(fmt.Sprintf("Reasoning: %s", msg))
}

func (m *ChatPage) renderAgentTextMessage(msg string) string {
	agentMsg := fmt.Sprintf("%s %s", lipgloss.NewStyle().Foreground(lipgloss.Color("#55ff66")).Render("Agent:"), msg)
	return agentMsg
}

func (m *ChatPage) renderErrorMessage(msg string) string {
	errMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render(fmt.Sprintf("Error: %s", msg))
	return errMsg
}

func (m *ChatPage) renderSeparator() string {
	separator := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render(strings.Repeat("─", m.viewport.Width()))
	return separator
}

func (m *ChatPage) close() {
	if m.currentAgent != nil {
		m.agentPool.Put(m.currentAgent)
	}
}

func (m *ChatPage) prepare() {
	if m.currentAgent == nil {
		a := m.agentPool.GetOrNew()
		m.currentAgent = a
	}
}
func (m *ChatPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.textarea.Height() - 2)
		m.viewport.GotoBottom()

	case StreamMsg:

		switch msg.ChunkType {
		case "done":
			m.messages = append(m.messages, m.renderSeparator())
		case "error":
			m.messages = append(m.messages, m.renderErrorMessage(msg.Content))
		case "text":
			if len(m.messages) > 0 && m.currentMessageChunkType == "text" {
				m.textMessage += msg.Content
				m.messages[len(m.messages)-1] = m.renderAgentTextMessage(m.textMessage)
			} else {
				m.textMessage = msg.Content
				m.messages = append(m.messages, m.renderAgentTextMessage(m.textMessage))
			}
		case "reasoning":
			if len(m.messages) > 0 && m.currentMessageChunkType == "reasoning" {
				m.reasoningMessage += msg.Content
				m.messages[len(m.messages)-1] = m.renderReasoningMessage(m.reasoningMessage)
			} else {
				m.reasoningMessage = msg.Content
				m.messages = append(m.messages, m.renderReasoningMessage(m.reasoningMessage))
			}
		case "tool_call":
			toolMsg := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#605959")).Render(fmt.Sprintf("-> Using tool: %s", msg.Content))
			m.messages = append(m.messages, toolMsg)
		}

		m.currentMessageChunkType = msg.ChunkType
		content := strings.Join(m.messages, "\n")
		wrapped := lipgloss.NewStyle().Width(m.viewport.Width()).Render(content)
		m.viewport.SetContent(wrapped)
		m.viewport.GotoBottom()
		cmds = append(cmds, waitForMessage(m.currentStream))
		return m, tea.Batch(cmds...)

	case StreamDoneMsg:
		m.sessionMgr.SaveHistory(m.sessionID, m.currentAgent.History)
		m.textarea.Placeholder = "Type a message... (Press Enter to send)"
		cmd := m.textarea.Focus()
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.close()
			return m, func() tea.Msg {
				return router.ChangePageMsg{
					Target: "main",
				}
			}
		case "enter":
			v := strings.TrimSpace(m.textarea.Value())
			if v == "" {
				break
			}

			// Wrap user input
			userMsg := fmt.Sprintf("You: %s", v)
			m.messages = append(m.messages, userMsg)
			content := strings.Join(m.messages, "\n")
			wrapped := lipgloss.NewStyle().Width(m.viewport.Width()).Render(content)
			m.viewport.SetContent(wrapped)
			m.viewport.GotoBottom()
			m.textarea.Reset()

			switch v {
			case "/exit":
				m.close()
				return m, func() tea.Msg {
					return router.ChangePageMsg{
						Target: "main",
					}
				}
			case "/settings":
				m.close()
				return m, func() tea.Msg {
					return router.ChangePageMsg{
						Target: "settings",
					}
				}
			}

			m.textarea.Blur()
			m.textarea.Placeholder = "Waiting for response..."

			m.prepare()
			m.currentAgent.History = m.sessionMgr.GetHistory(m.sessionID)
			m.currentStream = m.currentAgent.Invoke(context.Background(), v)
			m.currentMessageChunkType = "unknown"
			cmds = append(cmds, waitForMessage(m.currentStream))
			return m, tea.Batch(cmds...)
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

package chat

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

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
			separator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555555")).
				Render("────────────────────")
			m.messages = append(m.messages, separator)
		case "error":
			m.messages = append(m.messages, lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Render(fmt.Sprintf("Error: %s", msg.Content)))
		case "text":
			if len(m.messages) > 0 && m.currentMessageChunkType == "text" {
				m.messages[len(m.messages)-1] += msg.Content
			} else {
				m.messages = append(m.messages, fmt.Sprintf("%s %s", lipgloss.NewStyle().Foreground(lipgloss.Color("#55ff66")).Render("Agent:"), msg.Content))
			}
		case "reasoning":
			if len(m.messages) > 0 && m.currentMessageChunkType == "reasoning" {
				m.messages[len(m.messages)-1] += msg.Content
			} else {
				m.messages = append(m.messages, lipgloss.NewStyle().MarginLeft(10).Foreground(lipgloss.Color("#605959")).Render(fmt.Sprintf("Reasoning: %s", msg.Content)))
			}
		case "tool_call":
			m.messages = append(m.messages, lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#605959")).Render(fmt.Sprintf("-> Using tool: %s", msg.Content)))
		}

		m.currentMessageChunkType = msg.ChunkType
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		cmds = append(cmds, waitForMessage(m.currentStream))
		return m, tea.Batch(cmds...)

	case StreamDoneMsg:
		m.sessionMgr.SaveHistory(m.sessionID, m.currentAgent.History)
		m.agentPool.Put(m.currentAgent)
		m.currentAgent = nil

		m.textarea.Placeholder = "Type a message... (Press Enter to send)"
		cmd := m.textarea.Focus()
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.currentAgent != nil {
				break
			}
			v := strings.TrimSpace(m.textarea.Value())
			if v == "" {
				break
			}

			m.messages = append(m.messages, fmt.Sprintf("You: %s", v))
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.viewport.GotoBottom()
			m.textarea.Reset()

			m.textarea.Blur()
			m.textarea.Placeholder = "Waiting for response..."

			a := m.agentPool.GetOrNew()
			a.History = m.sessionMgr.GetHistory(m.sessionID)
			m.currentAgent = a
			m.currentStream = a.Invoke(context.Background(), v)
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

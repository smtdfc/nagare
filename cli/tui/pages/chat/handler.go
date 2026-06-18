package chat

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
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
		case "error":
			m.messages = append(m.messages, fmt.Sprintf("Error: %s\n", msg.Content))
		case "text":
			if len(m.messages) > 0 && m.currentMessageChunkType == "text" {
				m.messages[len(m.messages)-1] += msg.Content
			} else {
				m.messages = append(m.messages, fmt.Sprintf("Agent: %s", msg.Content))
			}

		case "tool_call":
			m.messages = append(m.messages, fmt.Sprintf("Using tool : %s\n", msg.Content))
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

			m.messages = append(m.messages, fmt.Sprintf("You: %s \n", v))
			m.viewport.SetContent(strings.Join(m.messages, "\n\n"))
			m.viewport.GotoBottom()
			m.textarea.Reset()

			m.textarea.Blur()
			m.textarea.Placeholder = "Waiting for response..."

			a := m.agentPool.GetOrNew()
			a.History = m.sessionMgr.GetHistory(m.sessionID)
			m.currentAgent = a
			m.currentStream = a.Invoke(context.Background(), v)

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

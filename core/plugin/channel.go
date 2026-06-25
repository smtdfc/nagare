package plugin

import (
	"context"

	"sync"

	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/messages"
	"github.com/smtdfc/nagare/plugin-sdk/host"
	"github.com/smtdfc/nagare/plugin-sdk/shared"
)

type ChatChannel struct {
	Id string

	Agent      *agent.Agent
	SessionMgr *agent.SessionManager
	SessionID  string
	CleanUp    func()
}

type ChatChannelManager struct {
	mu           sync.Mutex
	Host         *host.Host
	ChatChannels map[string]*ChatChannel
}

func (m *ChatChannelManager) Register(channel *ChatChannel) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ChatChannels[channel.Id] = channel
}

func (m *ChatChannelManager) Handle(payload *shared.HandleChatMessagePayload, pluginID string) {
	m.mu.Lock()
	channel := m.ChatChannels[payload.Channel]
	m.mu.Unlock()

	message := payload.Message
	channel.Agent.History = channel.SessionMgr.GetHistory(channel.SessionID)
	responseChannel := channel.Agent.Invoke(context.Background(), message)

	var fullResponse string

	for chunk := range responseChannel {
		switch msg := chunk.(type) {
		case *messages.TextMessage:
			fullResponse += msg.Content

		case *messages.StreamErrorMessage:
			fullResponse = "Oop! Nagare 😭: " + msg.Cause

		case *messages.ResponseFailedMessage:
			fullResponse = "Hummmm 🤡, Your LLM provider said: " + msg.Message

		case *messages.ToolCallMessage:
			m.Host.Send(pluginID, shared.HANDLE_CHAT_MESSAGE, shared.HandleChatMessagePayload{
				Channel: payload.Channel,
				Message: "🤖  Using " + msg.FunctionName,
			})
		}
	}

	m.Host.Send(pluginID, shared.HANDLE_CHAT_MESSAGE, shared.HandleChatMessagePayload{
		Channel: payload.Channel,
		Message: fullResponse,
	})
	channel.SessionMgr.SaveHistory(channel.SessionID, channel.Agent.History)
}

func NewChatChannelManager(host *host.Host) *ChatChannelManager {
	return &ChatChannelManager{
		mu:           sync.Mutex{},
		Host:         host,
		ChatChannels: map[string]*ChatChannel{},
	}
}

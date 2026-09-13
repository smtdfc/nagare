package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/smtdfc/nagare/plugin/client"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/message"
)

//go:embed metadata.json
var metadata string

const idleTimeout = 2 * time.Hour

type MessageBuffer struct {
	mu     sync.Mutex
	buffer strings.Builder
}

type TelegramPlugin struct {
	mu           sync.Mutex
	bot          *telego.Bot
	pluginClient *client.PluginClient
	channels     map[string]string
	idleTimers   map[string]*time.Timer

	queues         map[string][]telego.Update
	isProcessing   map[string]bool
	messageBuffers sync.Map
}

func NewTelegramPlugin(bot *telego.Bot, pluginClient *client.PluginClient) *TelegramPlugin {
	return &TelegramPlugin{
		bot:          bot,
		pluginClient: pluginClient,
		channels:     make(map[string]string),
		idleTimers:   make(map[string]*time.Timer),
		queues:       make(map[string][]telego.Update),
		isProcessing: make(map[string]bool),
	}
}

func (tp *TelegramPlugin) Start(ctx context.Context) error {
	_, err := tp.bot.GetMe(ctx)
	if err != nil {
		return err
	}

	_, err = tp.pluginClient.LoadMetadata(metadata)
	if err != nil {
		return err
	}

	tp.pluginClient.OnReceivedChatMessage = tp.OnReceivedChatMessage
	return tp.pluginClient.Start(
		ctx,
		func() {
			err := tp.pluginClient.Handshake(ctx)
			if err != nil {
				return
			}

			updates, _ := tp.bot.UpdatesViaLongPolling(ctx, nil)
			for update := range updates {
				if update.Message != nil {
					tp.enqueueMessage(update)
				}
			}
		},
	)
}

func (tp *TelegramPlugin) enqueueMessage(update telego.Update) {
	chatID := strconv.FormatInt(update.Message.Chat.ID, 10)

	tp.pluginClient.Logger.Info("enqueueMessage", "chatID", chatID)
	tp.mu.Lock()
	if t, ok := tp.idleTimers[chatID]; ok {
		t.Stop()
	}
	tp.idleTimers[chatID] = time.AfterFunc(idleTimeout, func() {
		tp.mu.Lock()
		if sessionID, active := tp.channels[chatID]; active {
			delete(tp.channels, chatID)
			delete(tp.idleTimers, chatID)
			delete(tp.queues, chatID)
			delete(tp.isProcessing, chatID)
			tp.pluginClient.Logger.Info("release chat session due to 2 hours of inactivity.", "chatID", chatID, "sessionID", sessionID)
		}
		tp.mu.Unlock()
	})

	tp.queues[chatID] = append(tp.queues[chatID], update)

	if !tp.isProcessing[chatID] {
		tp.isProcessing[chatID] = true
		tp.pluginClient.Logger.Info("processing", "chatID", chatID)
		go tp.processQueue(context.Background(), chatID)
	}

	tp.mu.Unlock()
}

func (tp *TelegramPlugin) processQueue(ctx context.Context, chatID string) {
	for {
		tp.mu.Lock()
		q := tp.queues[chatID]
		if len(q) == 0 {
			tp.isProcessing[chatID] = false
			tp.mu.Unlock()
			return
		}

		update := q[0]
		tp.queues[chatID] = q[1:]
		tp.mu.Unlock()
		tp.handleSingleMessage(ctx, &update)
	}
}

func (tp *TelegramPlugin) handleSingleMessage(ctx context.Context, update *telego.Update) {
	chatID := strconv.FormatInt(update.Message.Chat.ID, 10)

	tp.mu.Lock()
	sessionID, exist := tp.channels[chatID]
	tp.mu.Unlock()

	if !exist {
		var err error
		sessionID, err = tp.pluginClient.PrepareChatSession(ctx, chatID)
		if err != nil {
			_, err = tp.bot.SendMessage(ctx,
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					"Oops! Error while preparing chat session",
				),
			)

			if err != nil {
				tp.pluginClient.Logger.Error("Error sending message to telegram", "error", err)
			}
			return
		}

		tp.mu.Lock()
		tp.channels[chatID] = sessionID
		tp.mu.Unlock()
	}

	err := tp.pluginClient.SendChatMessage(ctx, sessionID, update.Message.Text)
	if err != nil {
		_, err = tp.bot.SendMessage(ctx,
			tu.Message(
				tu.ID(update.Message.Chat.ID),
				"Oops! Error while preparing chat session",
			),
		)

		if err != nil {
			tp.pluginClient.Logger.Error("Error sending message to telegram", "error", err)
		}
		return
	}
}

func (tp *TelegramPlugin) OnReceivedChatMessage(sessionID, chunk string) {
	val, _ := tp.messageBuffers.LoadOrStore(sessionID, &strings.Builder{})
	buf := val.(*strings.Builder)

	tp.mu.Lock()
	var targetChatID string
	for cID, sID := range tp.channels {
		if sID == sessionID {
			targetChatID = cID
			break
		}
	}
	tp.mu.Unlock()

	if targetChatID == "" {
		return
	}

	var chatIntID int64
	_, _ = fmt.Sscanf(targetChatID, "%d", &chatIntID)

	msg, err := helpers.UnmarshalJson[message.AnyMessage](chunk)
	if err != nil {

		return
	}

	switch msg.Type {
	case message.TextMessageType:
		msg, err := helpers.UnmarshalJson[message.TextMessage](chunk)
		if err != nil {
			return
		}

		tp.mu.Lock()
		buf.WriteString(msg.Content)
		tp.mu.Unlock()

	case message.AgentErrorMessageType:
		msg, err := helpers.UnmarshalJson[message.AgentErrorMessage](chunk)
		if err != nil {
			return
		}
		_, err = tp.bot.SendMessage(
			context.Background(),
			tu.Message(
				tu.ID(chatIntID),
				msg.Error,
			),
		)
		if err != nil {
			tp.pluginClient.Logger.Error("Error sending message to telegram", "error", err)
		}

		tp.messageBuffers.Delete(sessionID)

	case message.AgentCompletedMessageType:
		msg, err := helpers.UnmarshalJson[message.AgentCompletedMessage](chunk)
		if err != nil {
			return
		}

		if msg.Success {
			_, err = tp.bot.SendMessage(
				context.Background(),
				tu.Message(
					tu.ID(chatIntID),
					buf.String(),
				),
			)
			if err != nil {
				tp.pluginClient.Logger.Error("Error sending message to telegram", "error", err)
			}
		}

		tp.isProcessing[targetChatID] = false
		tp.messageBuffers.Delete(sessionID)
	}
}

func main() {
	token := os.Getenv("NAGARE_TELEGRAM_TOKEN")
	ctx := context.Background()
	bot, err := telego.NewBot(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pluginClient := client.NewPlugin()
	telegramPlugin := NewTelegramPlugin(bot, pluginClient)

	err = telegramPlugin.Start(ctx)
	if err != nil {
		fmt.Println("Plugin start error :", err)
	}
}

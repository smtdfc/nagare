package host

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/entities"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	core_plugin "github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/core/plugin/manager"
	"github.com/smtdfc/nagare/core/session"
	session_mgr "github.com/smtdfc/nagare/core/session/manager"
	task_manager "github.com/smtdfc/nagare/core/task/manager"
	tool_manager "github.com/smtdfc/nagare/core/tool/manager"
	"github.com/smtdfc/nagare/core/tool/registry"
	"github.com/smtdfc/nagare/plugin/client"
	"github.com/smtdfc/nagare/plugin/metadata"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AddInput struct {
	A int `json:"a"`
	B int `json:"b"`
}

type AddOutput struct {
	Sum int `json:"sum"`
}

func TestPluginHostLifecycleAndToolCall(t *testing.T) {
	tempDir := t.TempDir()
	sockPath := filepath.Join(tempDir, "plugin.sock")

	baseLogger, err := logger.NewBaseLogger()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite in memory: %v", err)
	}
	err = db.AutoMigrate(&entities.Plugin{}, &entities.Session{}, &entities.Task{}, &entities.Message{})
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}

	pluginRepo := repositories.NewPluginRepository(db, baseLogger)
	pluginMapper := mappers.NewPluginMapper()
	pluginMgr := manager.NewPluginManager(pluginRepo, pluginMapper, baseLogger)

	eventBus := event_bus.NewEventBus(baseLogger)
	taskRepo := repositories.NewTaskRepository(db, baseLogger)
	taskMgr := task_manager.NewTaskManager(taskRepo, baseLogger)
	toolMgr := tool_manager.NewToolManager(baseLogger, taskMgr, eventBus)
	sessionRepo := repositories.NewSessionRepository(db, baseLogger)
	sessionMapper := mappers.NewSessionMapper()
	messageRepo := repositories.NewMessageRepository(db, baseLogger)
	messageMapper := mappers.NewMessageMapper()
	sessionMgr := session_mgr.NewSessionManager(baseLogger, sessionRepo, sessionMapper, messageRepo, messageMapper, pluginRepo)

	// Seed plugin in DB
	testPlugin := &core_plugin.Plugin{
		PackageName: "nagare.test.math",
		Name:        "Math Plugin",
		Author:      "Nagare",
		Features:    []core_plugin.Feature{core_plugin.PluginToolFeature, core_plugin.ChatFeature},
		Version:     "1.0.0",
		Bin:         "/bin/true",
		IsActive:    true,
	}
	pluginEntity := pluginMapper.ToEntity(testPlugin)
	err = db.Create(pluginEntity).Error
	if err != nil {
		t.Fatalf("failed to create plugin in db: %v", err)
	}

	connectCode := uuid.NewString()
	pluginMgr.SetConnectCode(testPlugin.PackageName, connectCode)

	pluginHost := NewPluginHost(pluginMgr, sessionMgr, toolMgr, eventBus, baseLogger)
	pluginHost.SetSocketPath(sockPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pluginHost.Start(ctx)
	if err != nil {
		t.Fatalf("failed to start plugin host: %v", err)
	}
	defer func() {
		_ = pluginHost.Stop(ctx)
	}()

	// Define plugin client
	pluginClient := client.NewPlugin()
	pluginClient.Metadata = &metadata.PluginMetadata{
		PackageName: testPlugin.PackageName,
		Name:        testPlugin.Name,
		Author:      testPlugin.Author,
		Version:     testPlugin.Version,
		Features:    []string{"chat", "plugin_tool"},
	}
	pluginClient.ConnectConfig.SocketPath = sockPath
	pluginClient.ConnectConfig.ConnectCode = connectCode

	addTool := client.DefineTool(
		"math:add",
		"Adds two numbers",
		func(ctx *context.Context, args *AddInput) (*AddOutput, error) {
			return &AddOutput{Sum: args.A + args.B}, nil
		},
		[]string{"math"},
	)

	startedChan := make(chan struct{})
	clientErrChan := make(chan error, 1)

	go func() {
		err := pluginClient.Start(ctx, func() {
			close(startedChan)
		})
		if err != nil {
			clientErrChan <- err
		}
	}()

	select {
	case <-startedChan:
	case err := <-clientErrChan:
		t.Fatalf("client start failed: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for client to connect")
	}

	// Handshake
	err = pluginClient.Handshake(ctx)
	if err != nil {
		t.Fatalf("handshake failed: %v", err)
	}

	// Register tool
	err = pluginClient.RegisterPluginTool(ctx, addTool)
	if err != nil {
		t.Fatalf("register plugin tool failed: %v", err)
	}

	// Verify tool in registry
	toolItem, exists := registry.Registry["math:add"]
	if !exists {
		t.Fatal("expected math:add to be in registry")
	}

	// Call tool from host side
	execCtx := &core_context.ExecuteContext{
		Context:   ctx,
		SessionID: "session-1",
	}
	result, err := toolItem.Tool.Execute(execCtx, `{"a": 15, "b": 27}`)
	if err != nil {
		t.Fatalf("tool execution failed: %v", err)
	}

	if result != `{"sum":42}` {
		t.Fatalf("expected sum 42, got %s", result)
	}
}

func TestPluginHostChat(t *testing.T) {
	tempDir := t.TempDir()
	sockPath := filepath.Join(tempDir, "plugin_chat.sock")

	baseLogger, err := logger.NewBaseLogger()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	_ = db.AutoMigrate(&entities.Plugin{}, &entities.Session{}, &entities.Task{}, &entities.Message{})

	pluginRepo := repositories.NewPluginRepository(db, baseLogger)
	pluginMapper := mappers.NewPluginMapper()
	pluginMgr := manager.NewPluginManager(pluginRepo, pluginMapper, baseLogger)

	eventBus := event_bus.NewEventBus(baseLogger)
	taskRepo := repositories.NewTaskRepository(db, baseLogger)
	taskMgr := task_manager.NewTaskManager(taskRepo, baseLogger)
	toolMgr := tool_manager.NewToolManager(baseLogger, taskMgr, eventBus)
	sessionRepo := repositories.NewSessionRepository(db, baseLogger)
	sessionMapper := mappers.NewSessionMapper()
	messageRepo := repositories.NewMessageRepository(db, baseLogger)
	messageMapper := mappers.NewMessageMapper()
	sessionMgr := session_mgr.NewSessionManager(baseLogger, sessionRepo, sessionMapper, messageRepo, messageMapper, pluginRepo)

	testPlugin := &core_plugin.Plugin{
		PackageName: "nagare.test.chat",
		Name:        "Chat Plugin",
		Author:      "Nagare",
		Features:    []core_plugin.Feature{core_plugin.ChatFeature},
		Version:     "1.0.0",
		Bin:         "/bin/true",
		IsActive:    true,
	}
	_ = db.Create(pluginMapper.ToEntity(testPlugin)).Error

	connectCode := uuid.NewString()
	pluginMgr.SetConnectCode(testPlugin.PackageName, connectCode)

	pluginHost := NewPluginHost(pluginMgr, sessionMgr, toolMgr, eventBus, baseLogger)
	pluginHost.SetSocketPath(sockPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pluginHost.Start(ctx)
	if err != nil {
		t.Fatalf("failed to start plugin host: %v", err)
	}
	defer func() { _ = pluginHost.Stop(ctx) }()

	pluginClient := client.NewPlugin()
	pluginClient.Metadata = &metadata.PluginMetadata{
		PackageName: testPlugin.PackageName,
		Name:        testPlugin.Name,
		Author:      testPlugin.Author,
		Version:     testPlugin.Version,
		Features:    []string{"chat"},
	}
	pluginClient.ConnectConfig.SocketPath = sockPath
	pluginClient.ConnectConfig.ConnectCode = connectCode

	chunkReceived := make(chan string, 1)
	pluginClient.OnReceivedChatMessage = func(sessionID string, channelID string, chunk string) {
		chunkReceived <- chunk
	}

	startedChan := make(chan struct{})
	go func() {
		_ = pluginClient.Start(ctx, func() { close(startedChan) })
	}()

	select {
	case <-startedChan:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for client to connect")
	}

	err = pluginClient.Handshake(ctx)
	if err != nil {
		t.Fatalf("handshake failed: %v", err)
	}

	// Prepare session
	sessionID, err := pluginClient.PrepareChatSession(ctx, "channel-test")
	if err != nil {
		t.Fatalf("prepare chat session failed: %v", err)
	}
	if sessionID == "" {
		t.Fatal("expected non-empty sessionID")
	}

	// Listen on event bus for SendEvent
	sendChan, unsubscribeSend := eventBus.Subscribe(event_bus.SendEvent)
	defer unsubscribeSend()

	// Send chat message
	err = pluginClient.SendChatMessage(ctx, sessionID, "Hello assistant")
	if err != nil {
		t.Fatalf("send chat message failed: %v", err)
	}

	select {
	case evt := <-sendChan:
		sendPayload := evt.(*event_bus.SendMessageEventPayload)
		if sendPayload.Text != "Hello assistant" {
			t.Fatalf("expected 'Hello assistant', got %s", sendPayload.Text)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for SendEvent")
	}

	var pluginID string
	pluginHost.mu.RLock()
	for id := range pluginHost.connections {
		pluginID = id
	}
	pluginHost.mu.RUnlock()

	// Host forwards chunk to plugin
	eventBus.Publish(ctx, event_bus.ChunkEvent, &event_bus.ChatChunkEventPayload{
		SessionID:        sessionID,
		ChannelID:        "channel-test",
		SessionOwnerID:   pluginID,
		SessionOwnerType: string(session.PLUGIN),
	})

	select {
	case chunk := <-chunkReceived:
		_ = chunk
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for chunk")
	}

	// Reset channel
	err = pluginClient.ResetChatChannel(ctx, "channel-test")
	if err != nil {
		t.Fatalf("reset chat channel failed: %v", err)
	}
}

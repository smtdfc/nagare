package external_module

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/pkgs/paths"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type customData struct {
	Count int    `json:"count"`
	Text  string `json:"text"`
}

func TestDynamicActionEvent(t *testing.T) {
	// Kiểm tra tạo event động
	assert.Equal(t, Event("module:invoke:search"), MakeInvokeEvent("search"))
	assert.Equal(t, Event("module:invoke"), MakeInvokeEvent(""))
	assert.Equal(t, Event("module:reply:search"), MakeReplyEvent("search"))
	assert.Equal(t, Event("module:reply"), MakeReplyEvent(""))

	// Kiểm tra parse
	prefix, action := ParseActionEvent("module:invoke:vector_search")
	assert.Equal(t, "module:invoke", prefix)
	assert.Equal(t, "vector_search", action)

	prefix, action = ParseActionEvent("module:reply:click")
	assert.Equal(t, "module:reply", prefix)
	assert.Equal(t, "click", action)

	prefix, action = ParseActionEvent("module:reply")
	assert.Equal(t, "module:reply", prefix)
	assert.Equal(t, "", action)
}

func TestPayload_BindData(t *testing.T) {
	// Giả lập data nhận được từ JSON unmarshal là map[string]any
	rawMap := map[string]any{
		"count": float64(42),
		"text":  "hello world",
	}

	replyPayload := ReplyPayload[any]{
		Action:    "search",
		RequestID: "req-123",
		Data:      rawMap,
	}

	var target customData
	err := replyPayload.BindData(&target)
	require.NoError(t, err)
	assert.Equal(t, 42, target.Count)
	assert.Equal(t, "hello world", target.Text)
	assert.Equal(t, "search", replyPayload.Action)

	// Test DecodePayloadData helper
	decoded, err := DecodePayloadData[customData](rawMap)
	require.NoError(t, err)
	assert.Equal(t, 42, decoded.Count)
	assert.Equal(t, "hello world", decoded.Text)

	// Test InvokePayload BindData
	invokePayload := InvokePayload[any]{
		Action:    "search",
		RequestID: "req-456",
		Data:      rawMap,
	}
	var invokeTarget customData
	err = invokePayload.BindData(&invokeTarget)
	require.NoError(t, err)
	assert.Equal(t, 42, invokeTarget.Count)
	assert.Equal(t, "hello world", invokeTarget.Text)
	assert.Equal(t, "search", invokePayload.Action)
}

func TestExternalModuleHost_RegisterChannelsAndForward(t *testing.T) {
	log, err := logger.NewBaseLogger()
	require.NoError(t, err)
	eb := event_bus.NewEventBus(log)
	mgr := NewExternalModuleManager(log)
	host := NewExternalModuleHost(eb, log, mgr)

	moduleID := "test-module"

	// Đăng ký cặp kênh
	invokeCh, replyCh, err := host.RegisterChannels(moduleID, 10)
	require.NoError(t, err)
	assert.NotNil(t, invokeCh)
	assert.NotNil(t, replyCh)

	// Đăng ký trùng lặp phải báo lỗi
	_, _, err = host.RegisterChannels(moduleID, 10)
	assert.Error(t, err)

	// Giả lập mock net.Conn để kiểm tra forward
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	host.mu.Lock()
	host.connections[moduleID] = &Connection{
		conn:     serverConn,
		moduleID: moduleID,
	}
	host.mu.Unlock()

	// Gửi invoke qua invoke channel kèm action động
	go func() {
		invokeCh <- InvokePayload[any]{
			Action:    "ping",
			RequestID: "req-1",
			Data:      map[string]any{"msg": "hi"},
		}
	}()

	// Đọc từ clientConn (giả lập external module nhận tin)
	buf := make([]byte, 1024)
	n, err := clientConn.Read(buf)
	require.NoError(t, err)
	assert.True(t, n > 0)
	// Kiểm tra chuỗi chứa event động "module:invoke:ping"
	assert.Contains(t, string(buf[:n]), "module:invoke:ping")

	// Test forward reply động "module:reply:ping" từ external module về replyCh
	host.dispatch(&Connection{moduleID: moduleID}, &Payload[any]{
		RequestID: "req-1",
		Event:     MakeReplyEvent("ping"),
		Data:      "pong",
	})

	select {
	case reply := <-replyCh:
		assert.Equal(t, "req-1", reply.RequestID)
		assert.Equal(t, "ping", reply.Action)
		assert.Equal(t, "pong", reply.Data)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for reply")
	}

	// Test UnregisterChannels
	host.UnregisterChannels(moduleID)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err = RequestReply[string, string](ctx, invokeCh, replyCh, "test_action", "test_data", "req-2")
	assert.Error(t, err)
}

func TestExternalModuleHost_DedicatedSocketListener(t *testing.T) {
	log, err := logger.NewBaseLogger()
	require.NoError(t, err)
	eb := event_bus.NewEventBus(log)
	mgr := NewExternalModuleManager(log)
	host := NewExternalModuleHost(eb, log, mgr)

	moduleID := "dedicated-test-module"
	err = host.StartModuleListener(moduleID)
	require.NoError(t, err)
	defer host.Stop()

	// Verify listener map
	host.mu.RLock()
	listener, ok := host.listeners[moduleID]
	host.mu.RUnlock()
	assert.True(t, ok)
	assert.NotNil(t, listener)

	// Verify path generated matches
	expectedPath := paths.GetModuleSocketPath(moduleID)
	assert.Equal(t, expectedPath, listener.Addr().String())
}

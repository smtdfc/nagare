package manager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/hil"
	"github.com/smtdfc/nagare/core/tool"
	pkgs_messages "github.com/smtdfc/nagare/pkgs/messages"
)

func TestToolCallWaitsForApprovalBeforeWriting(t *testing.T) {
	approvalMgr := hil.NewManager()
	toolMgr := &ToolManager{
		approvalMgr: approvalMgr,
		logger:      approvalTestLogger(),
	}

	channel := make(chan pkgs_messages.Message, 1)
	execCtx := &core_context.ExecuteContext{
		Context:        context.Background(),
		MessageChannel: channel,
		SessionID:      "session-1",
	}

	path := filepath.Join(t.TempDir(), "approved.txt")

	go func() {
		msg := <-channel
		request, ok := msg.(*pkgs_messages.ApprovalRequestMessage)
		if !ok {
			t.Errorf("expected ApprovalRequestMessage, got %T", msg)
			return
		}
		if request.SessionID != "session-1" {
			t.Errorf("expected session-1, got %s", request.SessionID)
		}
		if err := approvalMgr.Respond(request.RequestID, "session-1", true); err != nil {
			t.Errorf("respond: %v", err)
		}
	}()

	result := toolMgr.Call(execCtx, &tool.Call{
		CallID: "call-1",
		Name:   "write_file_tool",
		Args:   fmt.Sprintf(`{"path":%q,"content":"approved","confirm":true}`, path),
	})

	if !result.IsSuccess {
		t.Fatalf("expected tool to succeed after approval, got %q", result.Result)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to be written, got: %v", err)
	}
}

func TestToolCallDoesNotWriteWhenApprovalIsDenied(t *testing.T) {
	approvalMgr := hil.NewManager()
	toolMgr := &ToolManager{
		approvalMgr: approvalMgr,
		logger:      approvalTestLogger(),
	}

	channel := make(chan pkgs_messages.Message, 1)
	execCtx := &core_context.ExecuteContext{
		Context:        context.Background(),
		MessageChannel: channel,
		SessionID:      "session-1",
	}

	path := filepath.Join(t.TempDir(), "denied.txt")

	go func() {
		msg := <-channel
		request, ok := msg.(*pkgs_messages.ApprovalRequestMessage)
		if !ok {
			t.Errorf("expected ApprovalRequestMessage, got %T", msg)
			return
		}
		if err := approvalMgr.Respond(request.RequestID, "session-1", false); err != nil {
			t.Errorf("respond: %v", err)
		}
	}()

	result := toolMgr.Call(execCtx, &tool.Call{
		CallID: "call-2",
		Name:   "write_file_tool",
		Args:   fmt.Sprintf(`{"path":%q,"content":"denied","confirm":true}`, path),
	})

	if result.IsSuccess {
		t.Fatalf("expected tool to fail after denial")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected file to not exist after denial, got: %v", err)
	}
}

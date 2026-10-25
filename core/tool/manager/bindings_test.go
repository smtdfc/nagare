package manager

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/hil"
	"github.com/smtdfc/nagare/core/logger"
	pkgs_messages "github.com/smtdfc/nagare/pkgs/messages"
)

func approvalTestLogger() *logger.BaseLogger {
	return &logger.BaseLogger{
		Logger: *slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func TestRequestToolApprovalFailsClosedWithoutContext(t *testing.T) {
	bindings := &ToolBindings{
		approvalMgr: hil.NewManager(),
		logger:      approvalTestLogger(),
	}

	if bindings.RequestToolApproval(nil, "delete a file") {
		t.Fatalf("expected approval to fail with a nil context")
	}
}

func TestRequestToolApprovalFailsClosedWithoutChannel(t *testing.T) {
	bindings := &ToolBindings{
		approvalMgr: hil.NewManager(),
		logger:      approvalTestLogger(),
	}

	ctx := &core_context.ExecuteContext{Context: context.Background()}
	if bindings.RequestToolApproval(ctx, "delete a file") {
		t.Fatalf("expected approval to fail with a nil message channel")
	}
}

func TestRequestToolApprovalAbortsWhenContextIsCancelled(t *testing.T) {
	bindings := &ToolBindings{
		approvalMgr: hil.NewManager(),
		logger:      approvalTestLogger(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	execCtx := &core_context.ExecuteContext{
		Context:        ctx,
		MessageChannel: make(chan pkgs_messages.Message),
	}

	done := make(chan bool, 1)
	go func() {
		done <- bindings.RequestToolApproval(execCtx, "delete a file")
	}()

	select {
	case approved := <-done:
		if approved {
			t.Fatalf("expected cancelled request to be denied")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("approval request blocked after context cancellation")
	}
}

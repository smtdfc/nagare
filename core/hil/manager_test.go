package hil

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestManagerAcceptsApprovalForMatchingSession(t *testing.T) {
	mgr := NewManager()
	request := mgr.Request("session-1", "Allow agent to delete file hello.txt ?")

	result := make(chan struct {
		accepted bool
		err      error
	}, 1)
	go func() {
		accepted, err := mgr.Wait(context.Background(), request.RequestID)
		result <- struct {
			accepted bool
			err      error
		}{accepted, err}
	}()

	if err := mgr.Respond(request.RequestID, "session-1", true); err != nil {
		t.Fatalf("respond: %v", err)
	}

	got := <-result
	if got.err != nil || !got.accepted {
		t.Fatalf("expected accepted approval, got accepted=%v err=%v", got.accepted, got.err)
	}
}

func TestManagerRejectsResponseForAnotherSession(t *testing.T) {
	mgr := NewManager()
	request := mgr.Request("session-1", "Allow agent to delete file hello.txt ?")

	err := mgr.Respond(request.RequestID, "session-2", true)
	if !errors.Is(err, ErrApprovalSessionMismatch) {
		t.Fatalf("expected session mismatch, got %v", err)
	}
}

func TestManagerTimesOut(t *testing.T) {
	mgr := NewManager()
	mgr.timeout = time.Millisecond
	request := mgr.Request("session-1", "Allow agent to delete file hello.txt ?")

	_, err := mgr.Wait(context.Background(), request.RequestID)
	if !errors.Is(err, ErrApprovalTimedOut) {
		t.Fatalf("expected timeout, got %v", err)
	}
}

func TestManagerCancelRemovesPendingRequest(t *testing.T) {
	mgr := NewManager()
	request := mgr.Request("session-1", "Allow agent to delete file hello.txt ?")

	mgr.Cancel(request.RequestID)

	err := mgr.Respond(request.RequestID, "session-1", true)
	if !errors.Is(err, ErrApprovalNotFound) {
		t.Fatalf("expected cancelled request to be gone, got %v", err)
	}
}

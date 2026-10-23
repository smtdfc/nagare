package messages

import "testing"

func TestNewTaskTriggerMessage(t *testing.T) {
	var taskTriggerTestMsg = NewTaskTriggerMessage()
	if taskTriggerTestMsg == nil {
		t.Fatal("expected NewTaskTriggerMessage, got <nil>")
	}
}

func TestTaskTriggerMessage_GetMessageType(t *testing.T) {
	var taskTriggerTestMsg = NewTaskTriggerMessage()
	if taskTriggerTestMsg.GetMessageType() != TaskTriggerMessageType {
		t.Fatalf("expected TaskTriggerMessageType, got %v", taskTriggerTestMsg.GetMessageType())
	}
}

func TestTaskTriggerMessage_GetMessageID(t *testing.T) {
	var taskTriggerTestMsg = NewTaskTriggerMessage()
	if taskTriggerTestMsg.GetMessageID() == "" {
		t.Fatal("expected message ID, got <empty string")
	}
}

func TestTaskTriggerMessage_GetInvokeID(t *testing.T) {
	var taskTriggerTestMsg = NewTaskTriggerMessage()
	taskTriggerTestMsg.SetInvokeID("inv-test")
	if taskTriggerTestMsg.GetInvokeID() != "inv-test" {
		t.Fatalf("expected invoke ID 'inv-start', got %q", taskTriggerTestMsg.GetInvokeID())
	}
}

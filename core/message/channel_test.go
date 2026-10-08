package message

import (
	"testing"

	"github.com/smtdfc/nagare/pkgs/messages"
)

func TestChannelTypes(t *testing.T) {
	// Verify Channel underlying channel conversions
	ch := make(Channel, 1)
	raw := (chan messages.Message)(ch)
	var writeCh WriteOnlyChannel = raw
	var readCh ReadOnlyChannel = raw

	msg := messages.NewTextMessage(messages.USER, "test message")
	writeCh <- msg

	received := <-readCh
	if received == nil || received.GetMessageType() != messages.TextMessageType {
		t.Errorf("expected received TextMessageType, got %v", received)
	}
}

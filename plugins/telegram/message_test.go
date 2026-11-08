package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/mymmrac/telego"
)

func TestSplitMessage(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		limit    int
		expected []string
	}{
		{
			name:     "empty",
			text:     "",
			limit:    10,
			expected: nil,
		},
		{
			name:     "fits in one message",
			text:     "hello",
			limit:    10,
			expected: []string{"hello"},
		},
		{
			name:     "exact limit",
			text:     "0123456789",
			limit:    10,
			expected: []string{"0123456789"},
		},
		{
			name:     "hard split",
			text:     "abcdefghij",
			limit:    4,
			expected: []string{"abcd", "efgh", "ij"},
		},
		{
			name:     "prefers line breaks",
			text:     "hello\nworld\nagain",
			limit:    10,
			expected: []string{"hello\n", "world\n", "again"},
		},
		{
			name:     "prefers spaces",
			text:     "hello world again",
			limit:    10,
			expected: []string{"hello ", "world ", "again"},
		},
		{
			name:     "keeps surrogate pairs intact",
			text:     "😀😀😀",
			limit:    4,
			expected: []string{"😀😀", "😀"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitMessage(tt.text, tt.limit)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d chunks, got %d: %q", len(tt.expected), len(got), got)
			}

			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("chunk %d: expected %q, got %q", i, tt.expected[i], got[i])
				}
				if tt.limit > 0 && utf16Length(got[i]) > tt.limit {
					t.Errorf("chunk %d exceeds limit: %d > %d", i, utf16Length(got[i]), tt.limit)
				}
			}

			if strings.Join(got, "") != tt.text {
				t.Errorf("chunks do not preserve the original text: %q", got)
			}
		})
	}
}

func TestSendTextMessage_SplitsLongMessages(t *testing.T) {
	var (
		mu    sync.Mutex
		texts []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ChatID int64  `json:"chat_id"`
			Text   string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		mu.Lock()
		texts = append(texts, body.Text)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1,"date":1,"chat":{"id":123,"type":"private"}}}`))
	}))
	t.Cleanup(server.Close)

	bot, err := telego.NewBot("123:"+strings.Repeat("a", 35), telego.WithAPIServer(server.URL))
	if err != nil {
		t.Fatalf("failed to create bot: %v", err)
	}

	tp := NewTelegramPlugin(bot, nil)
	text := strings.Repeat("a", telegramMaxMessageLength+10)

	if err := tp.sendTextMessage(context.Background(), 123, text); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(texts) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(texts))
	}
	if strings.Join(texts, "") != text {
		t.Errorf("sent messages do not preserve the original text")
	}
	for i, sent := range texts {
		if utf16Length(sent) > telegramMaxMessageLength {
			t.Errorf("message %d exceeds the Telegram limit: %d", i, utf16Length(sent))
		}
	}
}

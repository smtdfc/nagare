package ipc

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/smtdfc/nagare/pkgs/paths"
)

func getTestSocketPath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\nagare-test`
	}
	return filepath.Join(paths.DataDir, "nagare-test.sock")
}

func TestListenAndDial(t *testing.T) {
	testSocketPath := getTestSocketPath()

	if runtime.GOOS != "windows" {
		os.Remove(testSocketPath)
		defer os.Remove(testSocketPath)
	}

	s, err := Listen(testSocketPath)
	if err != nil {
		t.Fatalf("Failed to listen on named pipe: %v", err)
	}
	defer s.Close()

	go func() {
		conn, err := s.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
	}()

	ctx := context.Background()
	d, err := Dial(ctx, testSocketPath)
	if err != nil {
		t.Fatalf("Failed to dial named pipe: %v", err)
	}
	defer d.Close()
}

func TestListenAndSendMessage(t *testing.T) {
	testSocketPath := getTestSocketPath()

	if runtime.GOOS != "windows" {
		os.Remove(testSocketPath)
		defer os.Remove(testSocketPath)
	}

	s, err := Listen(testSocketPath)
	if err != nil {
		t.Fatalf("Failed to listen on named pipe: %v", err)
	}
	defer s.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		conn, err := s.Accept()
		if err != nil {
			t.Errorf("Failed to accept connection: %v", err)
			return
		}
		defer conn.Close()

		msg, err := ReadMessage(conn)
		if err != nil {
			t.Errorf("Failed to read message: %v", err)
			return
		}

		expectedMessage := "Hello, IPC!"
		if string(msg) != expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", expectedMessage, string(msg))
		}
	}()

	ctx := context.Background()
	d, err := Dial(ctx, testSocketPath)
	if err != nil {
		t.Fatalf("Failed to dial named pipe: %v", err)
	}
	defer d.Close()

	message := []byte("Hello, IPC!")
	err = WriteMessage(d, message)
	if err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	wg.Wait()
}

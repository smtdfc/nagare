package ipc

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/smtdfc/nagare/pkgs/paths"
)

const MaxMessageSize = 64 * 1024 * 1024 // 64 MB max message size

var (
	ErrMessageTooLarge = errors.New("ipc message exceeds maximum size")
)

// GetDefaultSocketPath returns the socket path or named pipe name.
func GetDefaultSocketPath() string {
	if envPath := os.Getenv("NAGARE_PLUGIN_SOCKET_PATH"); envPath != "" {
		return envPath
	}
	if runtime.GOOS == "windows" {
		return `\\.\pipe\nagare-plugin`
	}
	return paths.PluginSocketPath
}

// WriteMessage writes a length-prefixed frame to w.
func WriteMessage(w io.Writer, data []byte) error {
	if len(data) > MaxMessageSize {
		return ErrMessageTooLarge
	}
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)
	_, err := w.Write(buf)
	return err
}

// ReadMessage reads a length-prefixed frame from r.
func ReadMessage(r io.Reader) ([]byte, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	msgLen := binary.BigEndian.Uint32(lenBuf[:])
	if msgLen > MaxMessageSize {
		return nil, ErrMessageTooLarge
	}
	buf := make([]byte, msgLen)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// SendJSON marshals v and writes it as a length-prefixed frame to w.
func SendJSON(w io.Writer, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	return WriteMessage(w, data)
}

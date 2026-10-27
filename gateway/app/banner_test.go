package app

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestPrintStartupMessage(t *testing.T) {
	var buf bytes.Buffer

	printStartupMessage(&buf, "localhost:9832", true)

	got := buf.String()
	for _, want := range []string{
		"nagare gateway",
		"http://localhost:9832",
		"development",
		strconv.Itoa(os.Getpid()),
	} {
		if !strings.Contains(got, want) {
			t.Errorf("startup message missing %q, got:\n%s", want, got)
		}
	}
}

func TestPrintStartupMessage_Production(t *testing.T) {
	var buf bytes.Buffer

	printStartupMessage(&buf, "0.0.0.0:8080", false)

	if !strings.Contains(buf.String(), "production") {
		t.Errorf("expected production mode, got:\n%s", buf.String())
	}
}

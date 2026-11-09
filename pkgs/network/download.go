package network

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

const DefaultDownloadTimeout = 120 * time.Second
const DefaultDownloadMaxBytes = 1 << 30
const DefaultUserAgent = "nagare-agent/1.0"

type DownloadOptions struct {
	Timeout        time.Duration
	MaxBytes       int64
	Headers        map[string]string
	UserAgent      string
	ExpectedSha256 string
}

type DownloadResult struct {
	Bytes    int64
	SHA256   string
	Verified bool
}

func Download(rawURL string, destination string, options DownloadOptions) (DownloadResult, error) {
	if options.Timeout <= 0 {
		options.Timeout = DefaultDownloadTimeout
	}
	if options.MaxBytes <= 0 {
		options.MaxBytes = DefaultDownloadMaxBytes
	}
	if options.UserAgent == "" {
		options.UserAgent = DefaultUserAgent
	}

	request, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("failed to build request: %w", err)
	}
	request.Header.Set("User-Agent", options.UserAgent)
	for name, value := range options.Headers {
		request.Header.Set(name, value)
	}

	resp, err := (&http.Client{Timeout: options.Timeout}).Do(request)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return DownloadResult{}, fmt.Errorf("unsuccessful response status '%s' for %s", resp.Status, rawURL)
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(destination), ".nagare-download-*")
	if err != nil {
		return DownloadResult{}, fmt.Errorf("failed to create temporary file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	hasher := sha256.New()
	var written int64
	buf := make([]byte, 64*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			hasher.Write(buf[:n])
			written += int64(n)
			if written > options.MaxBytes {
				tmpFile.Close()
				return DownloadResult{}, fmt.Errorf("download exceeds the %d byte limit of this helper", options.MaxBytes)
			}
			if _, writeErr := tmpFile.Write(buf[:n]); writeErr != nil {
				tmpFile.Close()
				return DownloadResult{}, fmt.Errorf("failed to write download: %w", writeErr)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			tmpFile.Close()
			return DownloadResult{}, fmt.Errorf("failed to read download stream: %w", readErr)
		}
	}
	if err := tmpFile.Close(); err != nil {
		return DownloadResult{}, fmt.Errorf("failed to finalize temporary file: %w", err)
	}

	digest := hex.EncodeToString(hasher.Sum(nil))
	if options.ExpectedSha256 != "" && !strings.EqualFold(options.ExpectedSha256, digest) {
		return DownloadResult{}, fmt.Errorf("checksum mismatch: expected %s, got %s (file discarded)", options.ExpectedSha256, digest)
	}

	if err := os.Rename(tmpPath, destination); err != nil {
		return DownloadResult{}, fmt.Errorf("failed to finalize download: %w", err)
	}

	return DownloadResult{
		Bytes:    written,
		SHA256:   digest,
		Verified: options.ExpectedSha256 != "",
	}, nil
}

func BodyIsText(data []byte) bool {
	return utf8.Valid(data)
}

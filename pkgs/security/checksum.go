package security

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func ComputeFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file %s: %w", filePath, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func VerifyFileSignature(binFile, sigFile string) error {
	sigData, err := os.ReadFile(sigFile)
	if err != nil {
		return fmt.Errorf("failed to read signature file %s: %w", sigFile, err)
	}

	expectedSig := strings.TrimSpace(string(sigData))
	if expectedSig == "" {
		return errors.New("signature file is empty")
	}

	file, err := os.Open(binFile)
	if err != nil {
		return fmt.Errorf("failed to open binary file %s: %w", binFile, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to compute binary hash: %w", err)
	}

	actualSig := hex.EncodeToString(hasher.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(strings.ToLower(expectedSig)), []byte(strings.ToLower(actualSig))) != 1 {
		return fmt.Errorf("signature verification failed: expected %s, got %s", expectedSig, actualSig)
	}

	return nil
}

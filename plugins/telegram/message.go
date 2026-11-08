package main

import (
	"strings"
	"unicode/utf16"
)

// telegramMaxMessageLength is the maximum number of UTF-16 code units Telegram
// accepts in a single text message.
const telegramMaxMessageLength = 4096

func splitMessage(text string, limit int) []string {
	if text == "" {
		return nil
	}
	if limit <= 0 {
		limit = telegramMaxMessageLength
	}

	total := utf16Length(text)
	if total <= limit {
		return []string{text}
	}

	chunks := make([]string, 0, total/limit+1)
	remaining := text
	remainingLength := total

	for remaining != "" {
		if remainingLength <= limit {
			chunks = append(chunks, remaining)
			break
		}

		cut := byteIndexAtUTF16(remaining, limit)
		if cut <= 0 {
			break
		}

		if index := strings.LastIndexAny(remaining[:cut], "\n "); index > 0 && index >= cut/2 {
			cut = index + 1
		}

		chunk := remaining[:cut]
		chunks = append(chunks, chunk)
		remainingLength -= utf16Length(chunk)
		remaining = remaining[cut:]
	}

	return chunks
}

func utf16Length(text string) int {
	length := 0
	for _, r := range text {
		length += utf16.RuneLen(r)
	}
	return length
}

func byteIndexAtUTF16(text string, limit int) int {
	length := 0
	for index, r := range text {
		size := utf16.RuneLen(r)
		if length+size > limit {
			return index
		}
		length += size
	}
	return len(text)
}

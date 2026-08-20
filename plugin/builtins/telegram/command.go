package main

import (
	"strings"
)

type ParsedCommand struct {
	Command string
	Args    []string
	Raw     string
}

func ParseCommand(text string) *ParsedCommand {
	text = strings.TrimSpace(text)
	if text == "" || text[0] != '/' {
		return nil
	}

	fields := strings.Fields(text)
	if len(fields) == 0 {
		return nil
	}

	rawCmd := fields[0][1:]
	if idx := strings.Index(rawCmd, "@"); idx != -1 {
		rawCmd = rawCmd[:idx]
	}

	var args []string
	if len(fields) > 1 {
		args = fields[1:]
	}

	return &ParsedCommand{
		Command: strings.ToLower(rawCmd),
		Args:    args,
		Raw:     text,
	}
}

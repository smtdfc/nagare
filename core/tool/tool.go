package tool

import (
	"strings"

	"github.com/smtdfc/nagare/core/domains"
)

func GetToolSummary(tools domains.ListTool) string {
	if len(tools) == 0 {
		return ""
	}

	var b strings.Builder

	for _, t := range tools {
		b.WriteString(t.GetName())
		b.WriteString("(")
		b.WriteString(t.GetArgumentsSchema())
		b.WriteString("): ")
		b.WriteString(t.GetDesc())
		b.WriteString("\n")
	}

	return b.String()
}

func MergeToolLists(list1, list2 []domains.Tool) []domains.Tool {
	toolMap := make(map[string]domains.Tool)
	for _, t := range list1 {
		toolMap[t.GetName()] = t
	}
	for _, t := range list2 {
		toolMap[t.GetName()] = t
	}

	merged := make([]domains.Tool, 0, len(toolMap))
	for _, t := range toolMap {
		merged = append(merged, t)
	}

	return merged
}

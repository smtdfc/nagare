package plugin

import (
	"strings"
)

type PluginFeature string

const (
	BASIC          PluginFeature = "basic"
	AGENT          PluginFeature = "agent"
	SESSION_MANAGE PluginFeature = "session_manage"
)

func ParseFeaturesFromString(raw string) []PluginFeature {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	replacer := strings.NewReplacer(",", " ", ";", " ")
	cleaned := replacer.Replace(raw)

	fields := strings.Fields(cleaned)
	var features []PluginFeature

	for _, field := range fields {
		feature := PluginFeature(strings.ToLower(strings.TrimSpace(field)))
		switch feature {
		case BASIC, AGENT, SESSION_MANAGE:
			features = append(features, feature)
		}
	}

	return features
}

type ListPluginFeature []PluginFeature

func (l ListPluginFeature) ToString() string {
	if len(l) == 0 {
		return ""
	}
	strs := make([]string, len(l))
	for i, feature := range l {
		strs[i] = string(feature)
	}
	return strings.Join(strs, ", ")
}

func (l ListPluginFeature) ToSlice() []string {
	if len(l) == 0 {
		return nil
	}

	strs := make([]string, len(l))
	for i, feature := range l {
		strs[i] = string(feature)
	}
	return strs
}

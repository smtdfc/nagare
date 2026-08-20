package plugin

type PluginAuthPayload struct {
	ID       string
	Features ListPluginFeature
}

// GetID implements [AuthPayload].
func (j *PluginAuthPayload) GetID() string {
	return j.ID
}

// GetKind implements [AuthPayload].
func (j *PluginAuthPayload) GetKind() string {
	return "plugin"
}

// GetRole implements [AuthPayload].
func (j *PluginAuthPayload) GetRole() string {
	return ""
}

// GetScopes implements [AuthPayload].
func (j *PluginAuthPayload) GetScopes() []string {
	return j.Features.ToSlice()
}

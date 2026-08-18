package handlers

type PluginHandler struct {
}

// @Injectable
func NewPluginHandler() *PluginHandler {
	return &PluginHandler{}
}

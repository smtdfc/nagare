package tool

type ToolRegistry struct {
	StaticTools  []string
	DynamicTools []string
	Tools        ToolMap
}

func (r *ToolRegistry) Register(tool Tool) {
	toolType := tool.GetType()
	switch toolType {
	case STATIC_TOOL:
		r.StaticTools = append(r.StaticTools, tool.GetName())
	case DYNAMIC_TOOL:
		r.DynamicTools = append(r.DynamicTools, tool.GetName())
	}

	r.Tools[tool.GetName()] = tool
}

func (r *ToolRegistry) GetByName(name string) (Tool, bool) {
	t, ok := r.Tools[name]
	return t, ok
}

func (r *ToolRegistry) GetStaticTool() ListTool {
	tools := make(ListTool, len(r.StaticTools))
	for i, tool := range r.StaticTools {
		if t, exists := r.Tools[tool]; exists {
			tools[i] = t
		}
	}

	return tools
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		StaticTools:  []string{},
		DynamicTools: []string{},
		Tools:        make(ToolMap),
	}
}

var GlobalToolRegistry = NewToolRegistry()

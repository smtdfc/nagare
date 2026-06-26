package tool

import "github.com/smtdfc/nagare/core/domains"

type ToolRegistry struct {
	StaticTools  []string
	DynamicTools []string
	Tools        domains.ToolMap
}

func (r *ToolRegistry) Register(tool domains.Tool) {
	toolType := tool.GetType()
	switch toolType {
	case domains.STATIC_TOOL:
		r.StaticTools = append(r.StaticTools, tool.GetName())
	case domains.DYNAMIC_TOOL:
		r.DynamicTools = append(r.DynamicTools, tool.GetName())
	}

	r.Tools[tool.GetName()] = tool
}

func (r *ToolRegistry) GetByName(name string) (domains.Tool, bool) {
	t, ok := r.Tools[name]
	return t, ok
}

func (r *ToolRegistry) GetStaticTool() domains.ListTool {
	tools := make(domains.ListTool, len(r.StaticTools))
	for i, tool := range r.StaticTools {
		if t, exists := r.Tools[tool]; exists {
			tools[i] = t
		}
	}

	return tools
}

func (r *ToolRegistry) GetToolByCategories(categories []string) domains.ListTool {
	tools := make(domains.ListTool, 0)

	for _, toolName := range r.DynamicTools {
		if t, exists := r.Tools[toolName]; exists {
			if t.HasCategory(categories) {
				tools = append(tools, t)
			}
		}
	}

	return tools
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		StaticTools:  []string{},
		DynamicTools: []string{},
		Tools:        make(domains.ToolMap),
	}
}

package domains

type ToolType int

const (
	STATIC_TOOL ToolType = iota
	DYNAMIC_TOOL
)

const (
	PC_TOOL   string = "PC_TOOL"
	FILE_TOOL string = "FILE_TOOL"
)

type ListCategory []string

var NO_GROUP = []string{}

type Tool interface {
	GetName() string
	Execute(AgentContext, string) (string, error)
	// ExecuteWithRawResult(ctx context.Context, argsRaw string) (any, error)
	GetArgumentsSchema() string
	GetDesc() string
	GetType() ToolType
	HasCategory([]string) bool
}

type ListTool []Tool
type ToolMap map[string]Tool
type ToolCallResult struct {
	Result string
	Tool   Tool
}

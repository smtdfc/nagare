package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type GetUserDirectoriesOutput struct {
	Directories system.UserDirectories `json:"directories"`
}

var GetUserDirectoriesTool = tool.DefineTool(
	"get_user_directories_tool",
	"Get common user directories such as home, downloads, desktop, documents, temp, config, and cache.",
	func(ctx *context.ExecuteContext, _ *struct{}, _ tool.Bindings) (*GetUserDirectoriesOutput, error) {
		directories, err := system.GetUserDirectories()
		if err != nil {
			return nil, err
		}

		return &GetUserDirectoriesOutput{Directories: directories}, nil
	},
)

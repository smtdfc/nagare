package declarations

import (
	"os"
	"runtime"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type GetOSInfoInput struct{}

type GetOSInfoOutput struct {
	OS       string `json:"os" jsonschema_description:"Operating system family reported by the Go runtime, such as linux, darwin, or windows."`
	Arch     string `json:"arch" jsonschema_description:"CPU architecture reported by the Go runtime, such as amd64 or arm64."`
	Hostname string `json:"hostname" jsonschema_description:"Hostname of the machine running Nagare."`
	CPUs     int    `json:"cpus" jsonschema_description:"Number of logical CPUs available to the process."`
}

var GetOSInfoTool = tool.DefineTool(
	"get_os_info_tool",
	"Get operating system information: OS family, CPU architecture, hostname, and logical CPU count.",
	func(ctx *context.ExecuteContext, args *GetOSInfoInput, _ tool.Bindings) (*GetOSInfoOutput, error) {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = ""
		}

		return &GetOSInfoOutput{
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Hostname: hostname,
			CPUs:     runtime.NumCPU(),
		}, nil
	},
	[]string{tool.SystemCategory},
)

package declarations

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/cross-platform/system"
)

type GetSystemInfoArgs struct {
}

var get_system_info = tool.DeclareTool(
	"get_system_info",
	"Retrieve detailed system information.",
	func(ctx domains.AgentContext, args GetSystemInfoArgs) (any, error) {
		sysInfo, err := system.GetSystemInfo()
		if err != nil {
			return nil, err
		}

		return map[string]any{
			"os":       sysInfo.HostInfo.OS,
			"platform": sysInfo.HostInfo.Platform,
			"hostname": sysInfo.HostInfo.Hostname,
			"cpu": map[string]any{
				"model": sysInfo.CpuInfo[0].ModelName,
				"cores": sysInfo.CpuInfo[0].Cores,
			},
			"memory": map[string]any{
				"total_gb": sysInfo.Memory.Total / 1024 / 1024 / 1024,
				"used_pct": sysInfo.Memory.UsedPercent,
			},
		}, nil
	},
	domains.STATIC_TOOL,
	domains.NO_GROUP,
)

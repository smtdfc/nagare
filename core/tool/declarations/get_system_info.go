package declarations

import (
	"context"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/smtdfc/nagare/core/tool"
)

type GetSystemInfoArgs struct {
}

var get_system_info = tool.DeclareTool(
	"get_system_info",
	"Retrieve detailed system information including CPU, Memory, and OS details.",
	func(ctx context.Context, args GetSystemInfoArgs) (any, error) {
		cpuInfo, _ := cpu.Info()

		vMem, _ := mem.VirtualMemory()

		hostInfo, _ := host.Info()

		return map[string]any{
			"os":       hostInfo.OS,
			"platform": hostInfo.Platform,
			"hostname": hostInfo.Hostname,
			"cpu": map[string]any{
				"model": cpuInfo[0].ModelName,
				"cores": cpuInfo[0].Cores,
			},
			"memory": map[string]any{
				"total_gb": vMem.Total / 1024 / 1024 / 1024,
				"used_pct": vMem.UsedPercent,
			},
		}, nil
	},
	tool.STATIC_TOOL,
	tool.NO_GROUP,
)

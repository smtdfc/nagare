package system

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	CpuInfo  []cpu.InfoStat
	Memory   *mem.VirtualMemoryStat
	HostInfo *host.InfoStat
}

func GetSystemInfo() (*SystemInfo, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	vMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	return &SystemInfo{
		CpuInfo:  cpuInfo,
		Memory:   vMem,
		HostInfo: hostInfo,
	}, nil
}

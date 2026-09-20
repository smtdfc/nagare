package system

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcessQuery struct {
	MinMemoryMB   float32
	MinCPUPercent float64
	NameQuery     string
	ExcludeName   string
	SortBy        string
	Limit         int
}

type ProcessItem struct {
	PID        int32
	Name       string
	MemoryMB   float32
	CPUPercent float64
}

type ProcessManager struct{}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

func (pm *ProcessManager) Kill(pid int32) (string, error) {
	if pid <= 0 {
		return "", fmt.Errorf("invalid PID provided: %d", pid)
	}

	p, err := process.NewProcess(pid)
	if err != nil {
		return "", fmt.Errorf("process with PID %d not found or inaccessible: %w", pid, err)
	}

	name, _ := p.Name()
	if err := p.Kill(); err != nil {
		return name, fmt.Errorf("failed to terminate process: %w", err)
	}

	return name, nil
}

func (pm *ProcessManager) Find(query *ProcessQuery) ([]ProcessItem, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to list processes: %w", err)
	}

	var procList []ProcessItem

	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			continue
		}

		if query != nil && query.NameQuery != "" {
			if !strings.Contains(strings.ToLower(name), strings.ToLower(query.NameQuery)) {
				continue
			}
		}

		if query != nil && query.ExcludeName != "" {
			if strings.Contains(strings.ToLower(name), strings.ToLower(query.ExcludeName)) {
				continue
			}
		}

		memInfo, err := p.MemoryInfo()
		var memMB float32 = 0
		if err == nil && memInfo != nil {
			memMB = float32(memInfo.RSS) / (1024 * 1024)
		}

		if query != nil && query.MinMemoryMB > 0 && memMB < query.MinMemoryMB {
			continue
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}

		if query != nil && query.MinCPUPercent > 0 && cpuPercent < query.MinCPUPercent {
			continue
		}

		procList = append(procList, ProcessItem{
			PID:        p.Pid,
			Name:       name,
			MemoryMB:   memMB,
			CPUPercent: cpuPercent,
		})
	}

	if query != nil && query.SortBy != "" {
		sort.Slice(procList, func(i, j int) bool {
			switch strings.ToLower(query.SortBy) {
			case "ram", "memory":
				return procList[i].MemoryMB > procList[j].MemoryMB
			case "cpu":
				return procList[i].CPUPercent > procList[j].CPUPercent
			default:
				return false
			}
		})
	}

	if query != nil && query.Limit > 0 && len(procList) > query.Limit {
		procList = procList[:query.Limit]
	}

	return procList, nil
}

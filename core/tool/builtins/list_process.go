package declarations

import (
	"context"
	"fmt"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkg/system"
)

type ListProcessInput struct {
	MinMemoryMB   float32 `json:"min_memory_mb,omitempty"`   // Filter processes consuming RAM >= this amount (MB)
	MinCPUPercent float64 `json:"min_cpu_percent,omitempty"` // Filter processes consuming CPU >= this percentage (%)
	NameQuery     string  `json:"name_query,omitempty"`      // Filter by name (contains this keyword)
	ExcludeName   string  `json:"exclude_name,omitempty"`    // Exclude processes containing this keyword in their name
	SortBy        string  `json:"sort_by,omitempty"`         // "ram" or "cpu"
	Limit         int     `json:"limit,omitempty"`           // Limit the number of results (e.g., Top 10)
}

type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	MemoryMB   float32 `json:"memory_mb"`
	CPUPercent float64 `json:"cpu_percent"`
}

type ListProcessOutput struct {
	Processes []ProcessInfo `json:"processes"`
	Total     int           `json:"total"`
}

var ListProcessTool = tool.DefineTool(
	"list_process_tool",
	"List and filter running processes by memory, CPU, name inclusion, or name exclusion with sorting and limits",
	func(ctx context.Context, args *ListProcessInput) (*ListProcessOutput, error) {
		procMgr := system.NewProcessManager()

		var query *system.ProcessQuery
		if args != nil {
			query = &system.ProcessQuery{
				MinMemoryMB:   args.MinMemoryMB,
				MinCPUPercent: args.MinCPUPercent,
				NameQuery:     args.NameQuery,
				ExcludeName:   args.ExcludeName,
				SortBy:        args.SortBy,
				Limit:         args.Limit,
			}
		}

		items, err := procMgr.Find(query)
		if err != nil {
			return nil, fmt.Errorf("failed to find processes: %w", err)
		}

		var procList []ProcessInfo
		for _, item := range items {
			procList = append(procList, ProcessInfo{
				PID:        item.PID,
				Name:       item.Name,
				MemoryMB:   item.MemoryMB,
				CPUPercent: item.CPUPercent,
			})
		}

		return &ListProcessOutput{
			Processes: procList,
			Total:     len(procList),
		}, nil
	},
)

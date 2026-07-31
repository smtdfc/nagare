package services

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/smtdfc/nagare/server/custom_errors"
	"github.com/smtdfc/nagare/shared/dto"
)

type HealthService struct{}

func (h *HealthService) CheckHealth() (*dto.CheckHealthResponse, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, custom_errors.NewServiceError(err.Error(), 500)
	}

	vMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, custom_errors.NewServiceError(err.Error(), 500)
	}

	uptime, err := host.Uptime()
	if err != nil {
		return nil, custom_errors.NewServiceError(err.Error(), 500)
	}

	resp := dto.CheckHealthResponse{
		Cpu:    cpuPercent[0],
		Memory: vMem.UsedPercent,
		Uptime: int(uptime),
		Status: "OK",
	}

	return &resp, nil
}

// @Injectable
func NewHealthService() *HealthService {
	return &HealthService{}
}

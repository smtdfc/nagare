package dto

type CheckHealthResponse struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Uptime int     `json:"uptime"`
	Status string  `json:"status"`
}

package internal

import (
	plugin_dtos "github.com/smtdfc/nagare/plugin/dtos"
)

const (
	AddVectorEvent plugin_dtos.Event = "internal:vector:add"
)

type AddVectorEventPayload struct {
	RequestID string    `json:"requestID"`
	Vector    []float32 `json:"vector"`
	IDs       []uint64  `json:"ids"`
}

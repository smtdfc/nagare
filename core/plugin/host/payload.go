package host

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/smtdfc/nagare/dtos/websocket"
)

func getPayloadData[T any](payload *websocket.Payload[any]) (*T, error) {
	if payload == nil || payload.Data == nil {
		return nil, errors.New("payload or data is nil")
	}

	bytesData, err := json.Marshal(payload.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload data: %w", err)
	}

	var result T
	if err := json.Unmarshal(bytesData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into target type %T: %w", (*T)(nil), err)
	}

	return &result, nil
}

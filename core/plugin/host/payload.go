package host

import (
	"encoding/json"
	"errors"
	"fmt"

	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func getPayloadData[T any](payload *contracts.Payload[any]) (*T, error) {
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

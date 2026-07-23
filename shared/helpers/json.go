package helpers

import (
	"encoding/json"
	"fmt"
)

func ConvertToBytes(raw any) ([]byte, error) {
	switch v := raw.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T, expected string or []byte", raw)
	}
}

func FromJson[T any](raw any) (*T, error) {
	bytesData, err := ConvertToBytes(raw)
	if err != nil {
		return nil, err
	}

	var obj T
	if err := json.Unmarshal(bytesData, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

func MapObjectFromJson[T any](raw any, obj *T) (*T, error) {
	bytesData, err := ConvertToBytes(raw)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytesData, obj)
	return obj, err
}

func MapObjectToJson(obj any) (string, error) {
	bytesData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytesData), nil
}

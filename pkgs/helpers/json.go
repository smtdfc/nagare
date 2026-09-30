package helpers

import "encoding/json"

func UnmarshalJson[T any](raw string) (*T, error) {
	var data T
	err := json.Unmarshal([]byte(raw), &data)
	if err != nil {
		return &data, err
	}

	return &data, nil
}

func MarshalJson[T any](data T) (string, error) {
	raw, err := json.Marshal(&data)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

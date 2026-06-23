package utils

import (
	"encoding/json"
	"os"
)

func ReadJSON[T any](filePath string, v *T) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}

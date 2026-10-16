package host

import (
	"testing"

	plugin_dtos "github.com/smtdfc/nagare/plugin/dtos"
)

type sampleData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGetPayloadData(t *testing.T) {
	// Verify getPayloadData on valid payload
	payload := &plugin_dtos.Payload[any]{
		Data: map[string]any{
			"name": "Nagare",
			"age":  3,
		},
	}

	data, err := getPayloadData[sampleData](payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Name != "Nagare" || data.Age != 3 {
		t.Errorf("expected Nagare, 3; got %+v", data)
	}

	// Verify nil payload error
	_, err = getPayloadData[sampleData](nil)
	if err == nil {
		t.Errorf("expected error for nil payload")
	}

	// Verify nil data error
	nilDataPayload := &plugin_dtos.Payload[any]{Data: nil}
	_, err = getPayloadData[sampleData](nilDataPayload)
	if err == nil {
		t.Errorf("expected error for nil payload data")
	}
}

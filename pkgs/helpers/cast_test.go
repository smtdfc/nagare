package helpers

import (
	"testing"
)

func TestSafeCast(t *testing.T) {
	// Test successful type cast with string
	var strVal interface{} = "hello nagare"
	ok, val := SafeCast[string](strVal)
	if !ok || val != "hello nagare" {
		t.Fatalf("expected successful cast to string, got ok=%v, val=%v", ok, val)
	}

	// Test failed type cast with mismatched type
	ok, intVal := SafeCast[int](strVal)
	if ok || intVal != 0 {
		t.Fatalf("expected failed cast to int, got ok=%v, val=%v", ok, intVal)
	}

	// Test cast with custom struct
	type testStruct struct {
		Name string
		Age  int
	}
	var structVal interface{} = testStruct{Name: "Nagare", Age: 1}
	ok, castedStruct := SafeCast[testStruct](structVal)
	if !ok || castedStruct.Name != "Nagare" || castedStruct.Age != 1 {
		t.Fatalf("expected successful struct cast, got ok=%v, val=%+v", ok, castedStruct)
	}

	// Test cast with nil interface
	var nilVal interface{} = nil
	ok, nilCast := SafeCast[string](nilVal)
	if ok || nilCast != "" {
		t.Fatalf("expected failed cast for nil value, got ok=%v, val=%v", ok, nilCast)
	}
}

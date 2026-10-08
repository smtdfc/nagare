package helpers

import (
	"strconv"
	"testing"
)

func TestSliceMap(t *testing.T) {
	// Test normal transformation
	ints := []int{1, 2, 3, 4}
	res := SliceMap(ints, func(val int) *string {
		str := strconv.Itoa(val * 10)
		return &str
	})

	if len(res) != 4 {
		t.Fatalf("expected 4 mapped items, got %d", len(res))
	}
	expected := []string{"10", "20", "30", "40"}
	for i, v := range res {
		if v == nil || *v != expected[i] {
			t.Fatalf("index %d expected %s, got %v", i, expected[i], v)
		}
	}

	// Test filtering when callback returns nil
	mixed := []int{1, 2, 3, 4, 5}
	onlyEven := SliceMap(mixed, func(val int) *int {
		if val%2 == 0 {
			resVal := val
			return &resVal
		}
		return nil
	})

	if len(onlyEven) != 2 {
		t.Fatalf("expected 2 items after nil filter, got %d", len(onlyEven))
	}
	if *onlyEven[0] != 2 || *onlyEven[1] != 4 {
		t.Fatalf("expected [2, 4], got [%d, %d]", *onlyEven[0], *onlyEven[1])
	}

	// Test empty slice input
	empty := SliceMap([]string{}, func(s string) *string {
		return &s
	})
	if len(empty) != 0 {
		t.Fatalf("expected empty result for empty input, got %d", len(empty))
	}
}

package cmd

import (
	"reflect"
	"testing"
)

func TestNormalizeLegacyFlags(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "keeps regular args",
			in:   []string{"plugin", "build", "--sign", "key.pem"},
			want: []string{"plugin", "build", "--sign", "key.pem"},
		},
		{
			name: "expands single dash sign",
			in:   []string{"plugin", "build", "-sign", "key.pem"},
			want: []string{"plugin", "build", "--sign", "key.pem"},
		},
		{
			name: "expands single dash sign with value",
			in:   []string{"plugin", "build", "-sign=key.pem"},
			want: []string{"plugin", "build", "--sign=key.pem"},
		},
		{
			name: "leaves other shortcuts alone",
			in:   []string{"plugin", "build", "-s", "key.pem"},
			want: []string{"plugin", "build", "-s", "key.pem"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeLegacyFlags(tc.in); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("normalizeLegacyFlags(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

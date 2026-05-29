package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sum(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"1 + 1", 1, 1, 2},
		{"100 + 100", 100, 100, 200},
		{"0.1 + 0.2", 0.1, 0.2, 0.3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sum(tt.a, tt.b)
			require.Equal(t, tt.want, result)
		})
	}
}
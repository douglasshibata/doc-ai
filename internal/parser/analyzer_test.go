package parser

import (
	"testing"
)

func TestIsIntersecting(t *testing.T) {
	tests := []struct {
		start, end int
		changed    []int
		want       bool
	}{
		{10, 20, []int{15}, true},
		{10, 20, []int{5, 25}, false},
		{10, 20, []int{10}, true},
		{10, 20, []int{20}, true},
		{10, 20, []int{9, 11}, true},
	}

	for _, tt := range tests {
		if got := isIntersecting(tt.start, tt.end, tt.changed); got != tt.want {
			t.Errorf("isIntersecting(%d, %d, %v) = %v, want %v", tt.start, tt.end, tt.changed, got, tt.want)
		}
	}
}

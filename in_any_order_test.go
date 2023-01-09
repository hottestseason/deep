package deep_test

import (
	"testing"

	"github.com/hottestseason/deep"
)

func TestInAnyOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		want    []any
		args    any
		matches bool
	}{
		{
			name:    "1",
			want:    []any{1, 2, 3},
			args:    []int{1, 2, 3},
			matches: true,
		},
		{
			name:    "2",
			want:    []any{1, 2, 3},
			args:    []int{3, 1, 2},
			matches: true,
		},
		{
			name:    "3",
			want:    []any{},
			args:    []int{},
			matches: true,
		},
		{
			name:    "4",
			want:    nil,
			args:    []int{},
			matches: true,
		},
		{
			name:    "5",
			want:    []any{},
			args:    nil,
			matches: false,
		},
		{
			name:    "6",
			want:    []any{1, 2, 3},
			args:    []int{1, 4, 3},
			matches: false,
		},
		{
			name:    "7",
			want:    []any{1, 2, 3},
			args:    []int{3, 1, 1, 2},
			matches: false,
		},
		{
			name:    "8",
			want:    []any{},
			args:    []int{1, 2, 3},
			matches: false,
		},
		{
			name:    "8",
			want:    []any{1, 2, 3},
			args:    []int{},
			matches: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := deep.InAnyOrder(tt.want).Matches(tt.args); got != tt.matches {
				t.Errorf("deep.InAnyOrder(%v).Matches(%v) = %v, want %v", tt.want, tt.args, got, tt.matches)
			}
		})
	}
}

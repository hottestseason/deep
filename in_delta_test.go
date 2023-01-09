package deep_test

import (
	"testing"

	"github.com/hottestseason/deep"
	"golang.org/x/exp/constraints"
)

func TestInDelta(t *testing.T) {
	t.Parallel()

	// TODO: TDT
	testInDelta(t, 1, 2, -1.1, false)
	testInDelta(t, 1, 2, -1, true)
	testInDelta(t, 1, 2, -0.9, true)
	testInDelta(t, 1, 2, 2.9, true)
	testInDelta(t, 1, 2, 3, true)
	testInDelta(t, 1, 2, 3.1, false)
}

func testInDelta[T constraints.Integer | constraints.Float](t *testing.T, want T, delta float64, args any, matches bool) {
	if got := deep.InDelta(want, delta).Matches(args); got != matches {
		t.Errorf("deep.InDelta(%v, %v).Matches(%v) = %v, want %v", want, delta, args, got, matches)
	}
}

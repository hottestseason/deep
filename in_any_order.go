package deep

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slices"
)

type inAnyOrderMatcher[T any] struct {
	want []T
}

func InAnyOrder[T any](want []T) RecursiveMatcher {
	return inAnyOrderMatcher[T]{want}
}

func (m inAnyOrderMatcher[T]) Matches(x any) bool {
	return m.MatchesRecursively(x)
}

func (m inAnyOrderMatcher[T]) MatchesRecursively(x any, options ...cmp.Option) bool {
	vx := reflect.ValueOf(x)
	if vx.Kind() != reflect.Slice {
		return false
	}
	if len(m.want) != vx.Len() {
		return false
	}
	want := make([]T, len(m.want))
	copy(want, m.want)
	options = rejectReporters(options)
	for i := 0; i < vx.Len(); i++ {
		vxe := vx.Index(i).Interface()
		foundIdx := slices.IndexFunc(want, func(we T) bool { return cmp.Equal(we, vxe, options...) })
		if foundIdx == -1 {
			return false
		}
		want[foundIdx] = want[len(want)-1]
		want = want[:len(want)-1]
	}
	return len(want) == 0
}

func (m inAnyOrderMatcher[T]) String() string {
	return fmt.Sprintf("has the same elements as %v", m.want)
}

func rejectReporters(options []cmp.Option) []cmp.Option {
	ret := make([]cmp.Option, 0, len(options))
	for _, option := range options {
		if _, ok := option.(interface{ Report(cmp.Result) }); !ok {
			ret = append(ret, option)
		}
	}
	return ret
}

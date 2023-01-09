package deep

import (
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
)

// Matcher compatible to gomock.Matcher
type Matcher interface {
	Matches(x any) bool
	String() string
}

type RecursiveMatcher interface {
	Matcher
	MatchesRecursively(x any, options ...cmp.Option) bool
}

func Diff(want, got any) string {
	r := &diffReporter{}
	cmp.Equal(want, got, equate(cmp.Reporter(r))...)
	return strings.Join(r.diffs, "\n")
}

func equate(options ...cmp.Option) []cmp.Option {
	return append([]cmp.Option{
		cmp.FilterValues(func(x, y any) bool {
			_, okX := x.(Matcher)
			_, okY := y.(Matcher)
			return (okX && !okY) || (!okX && okY)
		}, cmp.Comparer(func(x, y any) bool {
			if x, okX := x.(RecursiveMatcher); okX {
				replaceLastMatcherPathStep(y, options)
				return x.MatchesRecursively(y, equate(options...)...)
			}
			if y, okY := y.(RecursiveMatcher); okY {
				replaceLastMatcherPathStep(x, options)
				return y.MatchesRecursively(x, equate(options...)...)
			}
			if x, okX := x.(Matcher); okX {
				replaceLastMatcherPathStep(y, options)
				return x.Matches(y)
			}
			if y, okY := y.(Matcher); okY {
				replaceLastMatcherPathStep(x, options)
				return y.Matches(x)
			}
			return false
		})),
		cmp.FilterValues(func(x, y any) bool { // Support []any to match []T
			vx := reflect.ValueOf(x)
			vy := reflect.ValueOf(y)
			return vx.Kind() == reflect.Slice && vy.Kind() == reflect.Slice && vx.Type() != vy.Type()
		}, cmp.Comparer(func(x, y any) bool {
			replaceLastMatcherPathStep(y, options)

			vx := reflect.ValueOf(x)
			vy := reflect.ValueOf(y)
			if vx.Len() != vy.Len() {
				return false
			}
			for i := 0; i < vx.Len(); i++ {
				vxe := vx.Index(i)
				vye := vy.Index(i)
				if !cmp.Equal(vxe.Interface(), vye.Interface(), equate(options...)...) {
					return false
				}
			}
			return true
		})),
	}, options...)
}

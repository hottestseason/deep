package deep

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
)

type superHashOfMatcher struct {
	want map[string]any
}

func SuperHashOf(want map[string]any) RecursiveMatcher {
	return superHashOfMatcher{want}
}

func (m superHashOfMatcher) Matches(x any) bool {
	return m.MatchesRecursively(x)
}

func (m superHashOfMatcher) MatchesRecursively(x any, options ...cmp.Option) bool {
	got, ok := x.(map[string]any)
	if !ok {
		return false
	}
	t := reflect.ValueOf(x).Type()
	popStep := pushStep(options, pathStep{
		typ:    t,
		vx:     reflect.ValueOf(m.want),
		vy:     reflect.ValueOf(x),
		string: ".(map[string]any)",
	})
	defer popStep()
	for k, v := range m.want { // TODO: for _, k := range value.SortKeys(append(vx.MapKeys(), vy.MapKeys()...)) {
		gotV, ok := got[k]
		if !ok {
			return false
		}
		popStep := pushStep(options, pathStep{
			typ:    t.Elem(),
			vx:     reflect.ValueOf(v),
			vy:     reflect.ValueOf(gotV),
			string: fmt.Sprintf("[%#v]", reflect.ValueOf(k)),
		})
		if !cmp.Equal(v, gotV, options...) {
			popStep()
			return false
		}
		popStep()
	}
	return true
}

func (m superHashOfMatcher) String() string {
	return fmt.Sprintf("is super hash of %v", m.want)
}

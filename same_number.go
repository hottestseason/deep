package deep

import (
	"fmt"
)

type sameNumberMatcher struct {
	want float64
}

func SameNumber(want any) Matcher {
	wantFloat, ok := toFloat(want)
	if !ok {
		panic("want must be number")
	}
	return sameNumberMatcher{wantFloat}
}

func (m sameNumberMatcher) Matches(x any) bool {
	xFloat, okX := toFloat(x)
	return okX && xFloat == m.want
}

func (m sameNumberMatcher) String() string {
	return fmt.Sprintf("%v", m.want)
}

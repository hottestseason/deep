package deep

import (
	"fmt"
	"math"
	"reflect"
)

type inDeltaMatcher struct {
	want  float64
	delta float64
}

func InDelta(want any, delta float64) Matcher {
	wantFloat, ok := toFloat(want)
	if !ok {
		panic("want must be number")
	}
	if delta < 0 {
		panic("delta must be positive number")
	}
	return inDeltaMatcher{wantFloat, delta}
}

func (m inDeltaMatcher) Matches(x any) bool {
	xFloat, okX := toFloat(x)
	return okX && math.Abs(xFloat-m.want) <= m.delta
}

func (m inDeltaMatcher) String() string {
	return fmt.Sprintf("is a number between %v and %v", m.want-m.delta, m.want+m.delta)
}

func toFloat(v interface{}) (float64, bool) {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(value.Uint()), true
	case reflect.Float32, reflect.Float64:
		return value.Float(), true
	default:
		return 0, false
	}
}

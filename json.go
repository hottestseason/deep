package deep

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
)

type jsonMatcher struct {
	want any
}

func JSON(want any) Matcher {
	return jsonMatcher{want}
}

func (m jsonMatcher) Matches(x any) bool {
	return m.MatchesRecursively(x)
}

func (m jsonMatcher) MatchesRecursively(x any, options ...cmp.Option) bool {
	got, ok := x.(string)
	if !ok {
		return false
	}
	var gotJSON any
	if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
		return false
	}
	return cmp.Equal(m.want, replaceFloatIntoSameNumber(gotJSON), options...)
}

func (m jsonMatcher) String() string {
	return fmt.Sprintf("is equal JSON of %v", m.want)
}

func replaceFloatIntoSameNumber(v any) any {
	vv := reflect.ValueOf(v)
	switch vv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return SameNumber(v)
	case reflect.Bool, reflect.String:
		return v
	case reflect.Map:
		for _, key := range vv.MapKeys() {
			vv.SetMapIndex(key, reflect.ValueOf(replaceFloatIntoSameNumber(vv.MapIndex(key).Interface())))
		}
		return vv.Interface()
	case reflect.Slice:
		for i := 0; i < vv.Len(); i++ {
			vv.Index(i).Set(reflect.ValueOf(replaceFloatIntoSameNumber(vv.Index(i).Interface())))
		}
		return vv.Interface()
	default:
		return v
	}
}

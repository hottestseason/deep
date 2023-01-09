package deep

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type diffReporter struct {
	path         cmp.Path
	values       [][2]reflect.Value
	lastDiffPath string
	diffs        []string
}

func (r *diffReporter) PushStep(ps cmp.PathStep) {
	if ps, ok := ps.(nextPathStep); ok {
		last := r.path.Last()
		if last.String() == "{any}" {
			r.path = append(r.path[:len(r.path)-1], ps)
		}
		return
	}
	r.path = append(r.path, ps)
	vx, vy := ps.Values()
	r.values = append(r.values, [2]reflect.Value{vx, vy})
}

func (r *diffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		path := r.path.GoString()
		if strings.HasPrefix(r.lastDiffPath, path) {
			return
		}
		lastValues := r.values[len(r.values)-1]
		printValue := func(v reflect.Value) string {
			if !v.IsValid() {
				return "<invalid reflect.Value>"
			}
			if v.Kind() == reflect.Interface && v.Elem().Kind() == reflect.Slice {
				return fmt.Sprintf("%+v (%T len:%d)", v.Interface(), v.Interface(), v.Elem().Len())
			}
			return fmt.Sprintf("%+v (%T)", v.Interface(), v.Interface())
		}
		r.diffs = append(r.diffs, fmt.Sprintf("%s:\n\t-: %s\n\t+: %s\n", path, printValue(lastValues[0]), printValue(lastValues[1])))
		r.lastDiffPath = path
	}
}

func (r *diffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
	r.values = r.values[:len(r.values)-1]
}

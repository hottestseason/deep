package deep

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
)

type pathStep struct {
	typ    reflect.Type
	vx, vy reflect.Value
	string string
}

func (ps pathStep) Type() reflect.Type             { return ps.typ }
func (ps pathStep) Values() (vx, vy reflect.Value) { return ps.vx, ps.vy }
func (ps pathStep) String() string                 { return ps.string }

func pushStep(options []cmp.Option, step pathStep) func() {
	for _, option := range options {
		if option, ok := option.(interface{ PushStep(cmp.PathStep) }); ok {
			option.PushStep(step)
		}
	}
	return func() {
		for _, option := range options {
			if option, ok := option.(interface{ PopStep() }); ok {
				option.PopStep()
			}
		}
	}
}

func replaceLastMatcherPathStep(v any, options []cmp.Option) {
	r := &nextPathStepExtractor{}
	defer func() {
		recover()
		for _, option := range options {
			if option, ok := option.(interface {
				PushStep(cmp.PathStep)
			}); ok {
				option.PushStep(r.step)
			}
		}
	}()
	cmp.Equal(v, v, cmp.Reporter(r))
}

type nextPathStep struct{ cmp.PathStep }

type nextPathStepExtractor struct {
	step cmp.PathStep
}

func (r *nextPathStepExtractor) PushStep(ps cmp.PathStep) {
	r.step = nextPathStep{ps}
	panic("finish")
}

func (r *nextPathStepExtractor) Report(rs cmp.Result) {}
func (r *nextPathStepExtractor) PopStep()             {}

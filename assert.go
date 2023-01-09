package deep

type TestingT interface {
	Errorf(format string, args ...any)
}

func AssertEqual(t TestingT, want any, got any) {
	if t, ok := t.(interface{ Helper() }); ok {
		t.Helper()
	}

	if diff := Diff(want, got); diff != "" {
		t.Errorf("diff:\n%s\n-want: %+v\n+got : %+v", diff, want, got)
	}
}

package try

import (
	"errors"
	"testing"
)

func TestMust_NoError(t *testing.T) {
	Must(nil)
}

func TestMust1_NoError(t *testing.T) {
	got := Must1(42, nil)
	if got != 42 {
		t.Fatalf("Must1 = %v, want 42", got)
	}
}

func TestMust2_NoError(t *testing.T) {
	a, b := Must2(1, "two", nil)
	if a != 1 || b != "two" {
		t.Fatalf("Must2 = (%v,%v), want (1,two)", a, b)
	}
}

func TestMust3_NoError(t *testing.T) {
	a, b, c := Must3(1, 2, 3, nil)
	if a != 1 || b != 2 || c != 3 {
		t.Fatalf("Must3 = (%v,%v,%v), want (1,2,3)", a, b, c)
	}
}

func TestMust4_NoError(t *testing.T) {
	a, b, c, d := Must4(1, 2, 3, 4, nil)
	if a != 1 || b != 2 || c != 3 || d != 4 {
		t.Fatalf("Must4 = (%v,%v,%v,%v), want (1,2,3,4)", a, b, c, d)
	}
}

func TestMust5_NoError(t *testing.T) {
	a, b, c, d, e := Must5(1, 2, 3, 4, 5, nil)
	if a != 1 || b != 2 || c != 3 || d != 4 || e != 5 {
		t.Fatalf("Must5 = (%v,%v,%v,%v,%v), want (1,2,3,4,5)", a, b, c, d, e)
	}
}

func TestMust_Error(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	Must(sentinelErr)
}

func TestMust1_Error(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	Must1(0, sentinelErr)
}

func TestMust2_Error(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	Must2(0, "", sentinelErr)
}

func TestMust3_Error(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	Must3(0, 0, 0, sentinelErr)
}

func TestMust4_Error(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	Must4(0, 0, 0, 0, sentinelErr)
}

func TestMust5_Error(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	Must5(0, 0, 0, 0, 0, sentinelErr)
}

func TestMust_PreservesErrorType(t *testing.T) {
	ce := &customErr{msg: "typed failure"}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic")
		}
		pe, ok := r.(panickedError)
		if !ok {
			t.Fatalf("panic value type = %T, want panickedError", r)
		}
		var target *customErr
		if !errors.As(pe, &target) {
			t.Fatalf("errors.As failed; unwrapped err = %T", pe.Unwrap())
		}
		if target != ce {
			t.Fatalf("target = %p, want %p", target, ce)
		}
	}()
	Must(ce)
}

func TestMust_OutsideScopePropagates(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	Must(sentinelErr)
}

func TestMust1_OutsideScopePropagates(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	_ = Must1(0, sentinelErr)
}

package try

import (
	"errors"
	"fmt"
	"testing"
)

// sentinelErr is a stable sentinel used across tests for errors.Is/As checks.
var sentinelErr = errors.New("sentinel failure")

// errFunc returns a typed error so we can verify errors.As round-trips
// through panickedError.
type customErr struct{ msg string }

func (e *customErr) Error() string { return e.msg }

// =============================================================================
// panickedError
// =============================================================================

func TestPanickedError_ErrorAndUnwrap(t *testing.T) {
	inner := fmt.Errorf("inner: %w", sentinelErr)
	pe := panickedError{err: inner}

	if got := pe.Error(); got != inner.Error() {
		t.Fatalf("Error() = %q, want %q", got, inner.Error())
	}
	if got := pe.Unwrap(); got != inner {
		t.Fatalf("Unwrap() = %v, want %v", got, inner)
	}
	// errors.As should see through panickedError to the wrapped error chain.
	var target *customErr
	ce := &customErr{msg: "typed"}
	pe2 := panickedError{err: ce}
	if !errors.As(pe2, &target) {
		t.Fatalf("errors.As(panickedError, *customErr) = false, want true")
	}
	if target != ce {
		t.Fatalf("target = %p, want %p", target, ce)
	}
}

// =============================================================================
// Must variants: nil error returns values without panic
// =============================================================================

func TestMust_NoError(t *testing.T) {
	Must(nil) // must not panic
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

// =============================================================================
// Must variants: non-nil error panics with panickedError wrapping it
// =============================================================================

func expectPanic(t *testing.T, want error) {
	t.Helper()
	r := recover()
	if r == nil {
		t.Fatalf("expected panic, got none")
	}
	pe, ok := r.(panickedError)
	if !ok {
		t.Fatalf("panic value = %#v, want panickedError", r)
	}
	if !errors.Is(pe, want) {
		t.Fatalf("errors.Is(panic.err, want) = false; got %v", pe.err)
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

// Verify Must preserves a concrete error type through panickedError.
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

// =============================================================================
// Scope: no panic returns nil error
// =============================================================================

func TestScope_NoPanic(t *testing.T) {
	called := false
	err := Scope(func() {
		called = true
	})
	if !called {
		t.Fatalf("f was not called")
	}
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
}

func TestScope1_NoPanic(t *testing.T) {
	out, err := Scope1(func() int { return 7 })
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if out != 7 {
		t.Fatalf("out = %v, want 7", out)
	}
}

func TestScope2_NoPanic(t *testing.T) {
	a, b, err := Scope2(func() (int, string) { return 1, "x" })
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if a != 1 || b != "x" {
		t.Fatalf("got (%v,%v), want (1,x)", a, b)
	}
}

func TestScope3_NoPanic(t *testing.T) {
	a, b, c, err := Scope3(func() (int, string, bool) { return 1, "x", true })
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if a != 1 || b != "x" || !c {
		t.Fatalf("got (%v,%v,%v), want (1,x,true)", a, b, c)
	}
}

func TestScope4_NoPanic(t *testing.T) {
	a, b, c, d, err := Scope4(func() (int, string, bool, float64) { return 1, "x", true, 2.5 })
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if a != 1 || b != "x" || !c || d != 2.5 {
		t.Fatalf("got (%v,%v,%v,%v), want (1,x,true,2.5)", a, b, c, d)
	}
}

func TestScope5_NoPanic(t *testing.T) {
	a, b, c, d, e, err := Scope5(func() (int, string, bool, float64, int) { return 1, "x", true, 2.5, 9 })
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if a != 1 || b != "x" || !c || d != 2.5 || e != 9 {
		t.Fatalf("got (%v,%v,%v,%v,%v), want (1,x,true,2.5,9)", a, b, c, d, e)
	}
}

// =============================================================================
// Scope: Must-induced panic is converted to err
// =============================================================================

func TestScope_MustPanicsConverted(t *testing.T) {
	err := Scope(func() {
		Must(sentinelErr)
	})
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("err = %v, want wrap of %v", err, sentinelErr)
	}
}

func TestScope1_MustPanicsConverted(t *testing.T) {
	out, err := Scope1(func() int {
		return Must1(0, sentinelErr)
	})
	if out != 0 {
		t.Fatalf("out = %v, want zero value", out)
	}
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("err = %v, want wrap of %v", err, sentinelErr)
	}
}

func TestScope2_MustPanicsConverted(t *testing.T) {
	a, b, err := Scope2(func() (int, string) {
		return Must2(0, "", sentinelErr)
	})
	if a != 0 || b != "" {
		t.Fatalf("got (%v,%v), want zero values", a, b)
	}
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("err = %v, want wrap of %v", err, sentinelErr)
	}
}

func TestScope3_MustPanicsConverted(t *testing.T) {
	a, b, c, err := Scope3(func() (int, string, bool) {
		return Must3(0, "", false, sentinelErr)
	})
	if a != 0 || b != "" || c != false {
		t.Fatalf("got (%v,%v,%v), want zero values", a, b, c)
	}
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("err = %v, want wrap of %v", err, sentinelErr)
	}
}

func TestScope4_MustPanicsConverted(t *testing.T) {
	a, b, c, d, err := Scope4(func() (int, string, bool, float64) {
		return Must4(0, "", false, 0.0, sentinelErr)
	})
	if a != 0 || b != "" || c != false || d != 0 {
		t.Fatalf("got (%v,%v,%v,%v), want zero values", a, b, c, d)
	}
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("err = %v, want wrap of %v", err, sentinelErr)
	}
}

func TestScope5_MustPanicsConverted(t *testing.T) {
	a, b, c, d, e, err := Scope5(func() (int, string, bool, float64, int) {
		return Must5(0, "", false, 0.0, 0, sentinelErr)
	})
	if a != 0 || b != "" || c != false || d != 0 || e != 0 {
		t.Fatalf("got (%v,%v,%v,%v,%v), want zero values", a, b, c, d, e)
	}
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("err = %v, want wrap of %v", err, sentinelErr)
	}
}

// =============================================================================
// Scope: non-Must panics are re-panicked (not swallowed)
// =============================================================================

func expectRepanic(t *testing.T, want any) {
	t.Helper()
	r := recover()
	if r == nil {
		t.Fatalf("expected Scope to re-panic, got none")
	}
	if fmt.Sprint(r) != fmt.Sprint(want) {
		t.Fatalf("re-panic value = %v (%T), want %v", r, r, want)
	}
}

func TestScope_RepanicsOnForeignPanic(t *testing.T) {
	defer expectRepanic(t, "foreign")
	_ = Scope(func() {
		panic("foreign")
	})
}

func TestScope_RepanicsOnNilDeref(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected re-panic from nil deref, got none")
		}
	}()
	_ = Scope(func() {
		var p *int
		_ = *p //nolint:staticcheck // intentional nil dereference
	})
}

func TestScope1_RepanicsOnForeignPanic(t *testing.T) {
	defer expectRepanic(t, "boom")
	_, _ = Scope1(func() int {
		panic("boom")
	})
}

func TestScope2_RepanicsOnForeignPanic(t *testing.T) {
	defer expectRepanic(t, "boom2")
	_, _, _ = Scope2(func() (int, int) {
		panic("boom2")
	})
}

func TestScope3_RepanicsOnForeignPanic(t *testing.T) {
	defer expectRepanic(t, "boom3")
	_, _, _, _ = Scope3(func() (int, int, int) {
		panic("boom3")
	})
}

func TestScope4_RepanicsOnForeignPanic(t *testing.T) {
	defer expectRepanic(t, "boom4")
	_, _, _, _, _ = Scope4(func() (int, int, int, int) {
		panic("boom4")
	})
}

func TestScope5_RepanicsOnForeignPanic(t *testing.T) {
	defer expectRepanic(t, "boom5")
	_, _, _, _, _, _ = Scope5(func() (int, int, int, int, int) {
		panic("boom5")
	})
}

// =============================================================================
// Integration: multi-step linear flow inside Scope
// =============================================================================

// stepN simulates a chain of fallible calls; each fails when its input is 0.
func step1(n int) (int, error) {
	if n == 0 {
		return 0, sentinelErr
	}
	return n + 1, nil
}
func step2(n int) (int, error) {
	if n == 0 {
		return 0, sentinelErr
	}
	return n * 2, nil
}

func TestScope_IntegrationShortCircuits(t *testing.T) {
	// step1(0) fails immediately; step2 must never be called.
	step2Called := false
	err := Scope(func() {
		a := Must1(step1(0)) // fails
		step2Called = true
		_ = Must1(step2(a))
	})
	if step2Called {
		t.Fatalf("step2 should not have been called after step1 failed")
	}
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("err = %v, want wrap of %v", err, sentinelErr)
	}
}

func TestScope_IntegrationHappyPath(t *testing.T) {
	err := Scope(func() {
		a := Must1(step1(5)) // 6
		b := Must1(step2(a)) // 12
		if b != 12 {
			panic(fmt.Sprintf("unexpected b = %d", b)) // foreign panic → would re-panic
		}
	})
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
}

func TestScope1_IntegrationReturnsValue(t *testing.T) {
	out, err := Scope1(func() int {
		a := Must1(step1(5))
		return Must1(step2(a))
	})
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if out != 12 {
		t.Fatalf("out = %v, want 12", out)
	}
}

func TestScope2_IntegrationReturnsValues(t *testing.T) {
	a, b, err := Scope2(func() (int, int) {
		x := Must1(step1(5))
		y := Must1(step2(x))
		return x, y
	})
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if a != 6 || b != 12 {
		t.Fatalf("got (%v,%v), want (6,12)", a, b)
	}
}

func TestScope3_IntegrationReturnsValues(t *testing.T) {
	a, b, c, err := Scope3(func() (int, int, int) {
		x := Must1(step1(5))
		y := Must1(step2(x))
		return x, y, x + y
	})
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if a != 6 || b != 12 || c != 18 {
		t.Fatalf("got (%v,%v,%v), want (6,12,18)", a, b, c)
	}
}

func TestScope4_IntegrationReturnsValues(t *testing.T) {
	a, b, c, d, err := Scope4(func() (int, int, int, int) {
		x := Must1(step1(5))
		y := Must1(step2(x))
		return x, y, x + y, x * y
	})
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if a != 6 || b != 12 || c != 18 || d != 72 {
		t.Fatalf("got (%v,%v,%v,%v), want (6,12,18,72)", a, b, c, d)
	}
}

func TestScope5_IntegrationReturnsValues(t *testing.T) {
	a, b, c, d, e, err := Scope5(func() (int, int, int, int, int) {
		x := Must1(step1(5))
		y := Must1(step2(x))
		return x, y, x + y, x * y, y - x
	})
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if a != 6 || b != 12 || c != 18 || d != 72 || e != 6 {
		t.Fatalf("got (%v,%v,%v,%v,%v), want (6,12,18,72,6)", a, b, c, d, e)
	}
}

// =============================================================================
// Must outside Scope: panic propagates (no recover)
// =============================================================================

func TestMust_OutsideScopePropagates(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	Must(sentinelErr)
}

func TestMust1_OutsideScopePropagates(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	_ = Must1(0, sentinelErr)
}

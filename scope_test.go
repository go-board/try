package try

import (
	"errors"
	"fmt"
	"testing"
)

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
	step2Called := false
	err := Scope(func() {
		a := Must1(step1(0))
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
		a := Must1(step1(5))
		b := Must1(step2(a))
		if b != 12 {
			panic(fmt.Sprintf("unexpected b = %d", b))
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

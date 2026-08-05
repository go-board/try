package try

import (
	"errors"
	"fmt"
	"testing"
)

func TestTry_MustHappyPath(t *testing.T) {
	Of(nil).Wrap("ignored").Wrapf("also ignored %d", 1).MapErr(func(error) error {
		t.Fatalf("MapErr should not be called for nil errors")
		return sentinelErr
	}).Must()

	if got := Of1(1, nil).Wrap("ignored").Must(); got != 1 {
		t.Fatalf("Of1(...).Must() = %v, want 1", got)
	}

	a, b := Of2(1, "two", nil).Wrapf("ignored %s", "context").Must()
	if a != 1 || b != "two" {
		t.Fatalf("Of2(...).Must() = (%v,%v), want (1,two)", a, b)
	}

	c, d, e := Of3(1, 2, 3, nil).Must()
	if c != 1 || d != 2 || e != 3 {
		t.Fatalf("Of3(...).Must() = (%v,%v,%v), want (1,2,3)", c, d, e)
	}

	f, g, h, i := Of4(1, 2, 3, 4, nil).Must()
	if f != 1 || g != 2 || h != 3 || i != 4 {
		t.Fatalf("Of4(...).Must() = (%v,%v,%v,%v), want (1,2,3,4)", f, g, h, i)
	}

	j, k, l, m, n := Of5(1, 2, 3, 4, 5, nil).Must()
	if j != 1 || k != 2 || l != 3 || m != 4 || n != 5 {
		t.Fatalf("Of5(...).Must() = (%v,%v,%v,%v,%v), want (1,2,3,4,5)", j, k, l, m, n)
	}
}

func TestTry_WrapPreservesErrorIdentity(t *testing.T) {
	err := Scope(func() {
		Of(sentinelErr).Wrap("open config").Must()
	})
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
	if got, want := err.Error(), "open config: sentinel failure"; got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}
}

func TestTry_WrapfPreservesErrorIdentity(t *testing.T) {
	err := Scope(func() {
		_ = Of1(0, sentinelErr).Wrapf("parse item %d", 42).Must()
	})
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
	if got, want := err.Error(), "parse item 42: sentinel failure"; got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}
}

func TestTry_MapErrTransformsError(t *testing.T) {
	mappedErr := errors.New("mapped")
	err := Scope(func() {
		_, _ = Of2(0, "", sentinelErr).MapErr(func(err error) error {
			return fmt.Errorf("%w: %w", mappedErr, err)
		}).Must()
	})
	if !errors.Is(err, mappedErr) {
		t.Fatalf("errors.Is(err, mappedErr) = false; got %v", err)
	}
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
}

func TestTry_MapErrNilKeepsOriginalError(t *testing.T) {
	err := Scope(func() {
		_, _, _ = Of3(0, 0, 0, sentinelErr).MapErr(func(error) error {
			return nil
		}).Must()
	})
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
}

func TestTry_MultipleValueWraps(t *testing.T) {
	err := Scope(func() {
		_, _, _, _ = Of4(0, 0, 0, 0, sentinelErr).Wrap("load tuple").Must()
	})
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}

	err = Scope(func() {
		_, _, _, _, _ = Of5(0, 0, 0, 0, 0, sentinelErr).Wrapf("load tuple %d", 5).Must()
	})
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
	if got, want := err.Error(), "load tuple 5: sentinel failure"; got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}
}

func TestTry_ResultHappyPath(t *testing.T) {
	if err := Of(nil).Result(); err != nil {
		t.Fatalf("Of(nil).Result() = %v, want nil", err)
	}

	a, err := Of1(1, nil).Result()
	if err != nil || a != 1 {
		t.Fatalf("Of1(...).Result() = (%v,%v), want (1,nil)", a, err)
	}

	b, c, err := Of2(1, "two", nil).Result()
	if err != nil || b != 1 || c != "two" {
		t.Fatalf("Of2(...).Result() = (%v,%v,%v), want (1,two,nil)", b, c, err)
	}

	d, e, f, err := Of3(1, 2, 3, nil).Result()
	if err != nil || d != 1 || e != 2 || f != 3 {
		t.Fatalf("Of3(...).Result() = (%v,%v,%v,%v), want (1,2,3,nil)", d, e, f, err)
	}

	g, h, i, j, err := Of4(1, 2, 3, 4, nil).Result()
	if err != nil || g != 1 || h != 2 || i != 3 || j != 4 {
		t.Fatalf("Of4(...).Result() = (%v,%v,%v,%v,%v), want (1,2,3,4,nil)", g, h, i, j, err)
	}

	k, l, m, n, o, err := Of5(1, 2, 3, 4, 5, nil).Result()
	if err != nil || k != 1 || l != 2 || m != 3 || n != 4 || o != 5 {
		t.Fatalf("Of5(...).Result() = (%v,%v,%v,%v,%v,%v), want (1,2,3,4,5,nil)", k, l, m, n, o, err)
	}
}

func TestTry_ResultReturnsWrappedError(t *testing.T) {
	err := Of(sentinelErr).Wrap("commit").Result()
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
	if got, want := err.Error(), "commit: sentinel failure"; got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}

	_, err = Of1(0, sentinelErr).Wrapf("load item %d", 7).Result()
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
	if got, want := err.Error(), "load item 7: sentinel failure"; got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}
}

func TestTry_ResultReturnsMappedError(t *testing.T) {
	mappedErr := errors.New("mapped")
	_, _, err := Of2(0, "", sentinelErr).MapErr(func(err error) error {
		return fmt.Errorf("%w: %w", mappedErr, err)
	}).Result()
	if !errors.Is(err, mappedErr) {
		t.Fatalf("errors.Is(err, mappedErr) = false; got %v", err)
	}
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
}

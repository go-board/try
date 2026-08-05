package try

import (
	"errors"
	"testing"
)

func TestAssert_NoPanicWhenTrue(t *testing.T) {
	err := Scope(func() {
		Assert(true, sentinelErr)
		Assertf(true, "failure: %w", sentinelErr)
	})
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
}

func TestAssert_FailureConvertedByScope(t *testing.T) {
	err := Scope(func() {
		Assert(false, sentinelErr)
	})
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
}

func TestAssert_NilErrorUsesDefault(t *testing.T) {
	err := Scope(func() {
		Assert(false, nil)
	})
	if !errors.Is(err, ErrAssertionFailed) {
		t.Fatalf("errors.Is(err, ErrAssertionFailed) = false; got %v", err)
	}
}

func TestAssertf_FailureConvertedByScope(t *testing.T) {
	err := Scope(func() {
		Assertf(false, "value %d: %w", 7, sentinelErr)
	})
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
	if got, want := err.Error(), "value 7: sentinel failure"; got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}
}

func TestAssert_OutsideScopePropagates(t *testing.T) {
	defer expectPanic(t, sentinelErr)
	Assert(false, sentinelErr)
}

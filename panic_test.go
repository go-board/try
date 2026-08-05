package try

import (
	"errors"
	"fmt"
	"testing"
)

func TestPanickedError_ErrorAndUnwrap(t *testing.T) {
	inner := fmt.Errorf("inner: %w", sentinelErr)
	pe := panickedError{err: inner}

	if got := pe.Error(); got != inner.Error() {
		t.Fatalf("Error() = %q, want %q", got, inner.Error())
	}
	if got := pe.Unwrap(); got != inner {
		t.Fatalf("Unwrap() = %v, want %v", got, inner)
	}

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

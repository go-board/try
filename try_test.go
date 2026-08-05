package try

import (
	"errors"
	"testing"
)

// sentinelErr is a stable sentinel used across tests for errors.Is/As checks.
var sentinelErr = errors.New("sentinel failure")

// customErr verifies that errors.As round-trips through panickedError.
type customErr struct{ msg string }

func (e *customErr) Error() string { return e.msg }

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

// Package try provides a flat, exception-style error-handling pattern built on
// top of Go's native error values.
//
// The idea: inside a [Scope] block, fallible calls are wrapped with [Must] /
// [Must1] .. [Must5] instead of the usual `if err != nil { return err }`
// ceremony. A Must failure raises an internal panic carrying the original
// error; the surrounding Scope recovers that panic and returns it as an
// ordinary error. Callers therefore write linear, panic-free business logic
// while still exposing errors as values at the function boundary.
//
// Only panics raised by this package's Must family are captured by Scope. Any
// other panic (nil dereference, out-of-range, third-party library panic) is
// re-panicked so genuine bugs are never silently swallowed.
//
// # Quick start
//
//	a, err := try.Scope1(func() int {
//	    x := try.Must1(step1(in))   // linear, no err checks
//	    y := try.Must1(step2(x))
//	    return try.Must1(step3(y))
//	})
//	// a, err behave like any (value, error) pair: errors.Is/As work, and the
//	// wrapped error identity is preserved through panickedError.
//
// # Concurrency
//
// Must is safe for concurrent use. Scope only captures panics from the same
// goroutine that invokes it: a goroutine started inside f that calls Must and
// panics will crash the program, which is the only safe default since a parent
// goroutine cannot reliably recover a child's panic.
package try

// panickedError is the sentinel type used to distinguish panics raised by this
// package's Must family (intentional control flow) from unrelated panics (real
// bugs that should propagate). It is unexported on purpose: callers interact
// only with the wrapped error, never with panickedError itself.
//
// Scope recovers a panic and type-asserts it to panickedError. The assertion
// fails for any foreign panic value, in which case Scope re-panics with the
// original value so the bug surfaces at its true origin.
type panickedError struct{ err error }

// Error implements error. The message is forwarded from the wrapped error so
// panickedError is transparent for logging and human-readable output.
func (e panickedError) Error() string { return e.err.Error() }

// Unwrap returns the original error, keeping panickedError compatible with
// [errors.Is] and [errors.As]. This means the concrete type and sentinel
// identity of the underlying error survive a Must -> Scope round-trip.
func (e panickedError) Unwrap() error { return e.err }

// checkError is the single chokepoint for the Must family: when err is non-nil
// it raises a panickedError wrapping err, otherwise it returns silently.
//
// Centralizing the panic here keeps the stack frame and behavior identical
// across Must / Must1 .. Must5, so Scope's recover logic has one invariant
// to match against rather than six.
func checkError(err error) {
	if err != nil {
		panic(panickedError{err: err})
	}
}

// Must panics if err is non-nil.
//
// It is intended for use inside a [Scope] block: the panic is captured by
// Scope and returned to the caller as an error. Used outside Scope the panic
// propagates and crashes the program, similar to [template.Must].
//
// Panics raised by Must wrap the original error in a panickedError, so the
// underlying value is still reachable via errors.Is / errors.As after Scope
// converts the panic back into an error.
func Must(err error) {
	checkError(err)
}

// Must1 returns v if err is nil; otherwise it panics with a panickedError
// wrapping err. See [Must] for the intended usage inside [Scope].
func Must1[A any](v A, err error) A {
	checkError(err)
	return v
}

// Must2 returns v1 and v2 if err is nil; otherwise it panics with a
// panickedError wrapping err. See [Must] for the intended usage inside [Scope].
func Must2[A, B any](v1 A, v2 B, err error) (A, B) {
	checkError(err)
	return v1, v2
}

// Must3 returns v1, v2 and v3 if err is nil; otherwise it panics with a
// panickedError wrapping err. See [Must] for the intended usage inside [Scope].
func Must3[A, B, C any](v1 A, v2 B, v3 C, err error) (A, B, C) {
	checkError(err)
	return v1, v2, v3
}

// Must4 returns v1 through v4 if err is nil; otherwise it panics with a
// panickedError wrapping err. See [Must] for the intended usage inside [Scope].
func Must4[A, B, C, D any](v1 A, v2 B, v3 C, v4 D, err error) (A, B, C, D) {
	checkError(err)
	return v1, v2, v3, v4
}

// Must5 returns v1 through v5 if err is nil; otherwise it panics with a
// panickedError wrapping err. See [Must] for the intended usage inside [Scope].
func Must5[A, B, C, D, E any](v1 A, v2 B, v3 C, v4 D, v5 E, err error) (A, B, C, D, E) {
	checkError(err)
	return v1, v2, v3, v4, v5
}

// Scope executes f and converts any [Must]-induced panic back into an error
// returned via err.
//
// Panics whose value is not a panickedError are re-panicked unchanged, so real
// bugs (nil dereference, out-of-range, third-party panics) propagate with
// their original stack and are never silently swallowed as errors.
//
// Only panics raised on the same goroutine that calls Scope are captured. If
// f spawns goroutines that themselves call Must, their panics will escape
// Scope and crash the program; a parent goroutine cannot reliably recover a
// child's panic.
//
// On the happy path (no panic) err is nil.
func Scope(f func()) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		pe, ok := r.(panickedError)
		if !ok {
			panic(r) // re-panic: not from our Must
		}
		err = pe.err
	}()
	f()
	return
}

// Scope1 is like [Scope] but f returns a single value, returned alongside the
// error. On a Must-induced panic, out is the zero value of A and err carries
// the original error. Foreign panics are re-panicked, as with Scope.
func Scope1[A any](f func() A) (out A, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		pe, ok := r.(panickedError)
		if !ok {
			panic(r) // re-panic: not from our Must
		}
		err = pe.err
	}()
	out = f()
	return
}

// Scope2 is like [Scope1] but f returns two values. On a Must-induced panic,
// both outputs are zero-valued and err carries the original error.
func Scope2[A, B any](f func() (A, B)) (out1 A, out2 B, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		pe, ok := r.(panickedError)
		if !ok {
			panic(r) // re-panic: not from our Must
		}
		err = pe.err
	}()
	out1, out2 = f()
	return
}

// Scope3 is like [Scope1] but f returns three values. On a Must-induced
// panic, all outputs are zero-valued and err carries the original error.
func Scope3[A, B, C any](f func() (A, B, C)) (out1 A, out2 B, out3 C, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		pe, ok := r.(panickedError)
		if !ok {
			panic(r) // re-panic: not from our Must
		}
		err = pe.err
	}()
	out1, out2, out3 = f()
	return
}

// Scope4 is like [Scope1] but f returns four values. On a Must-induced panic,
// all outputs are zero-valued and err carries the original error.
func Scope4[A, B, C, D any](f func() (A, B, C, D)) (out1 A, out2 B, out3 C, out4 D, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		pe, ok := r.(panickedError)
		if !ok {
			panic(r) // re-panic: not from our Must
		}
		err = pe.err
	}()
	out1, out2, out3, out4 = f()
	return
}

// Scope5 is like [Scope1] but f returns five values. On a Must-induced panic,
// all outputs are zero-valued and err carries the original error.
func Scope5[A, B, C, D, E any](f func() (A, B, C, D, E)) (out1 A, out2 B, out3 C, out4 D, out5 E, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		pe, ok := r.(panickedError)
		if !ok {
			panic(r) // re-panic: not from our Must
		}
		err = pe.err
	}()
	out1, out2, out3, out4, out5 = f()
	return
}

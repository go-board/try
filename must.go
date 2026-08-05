package try

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

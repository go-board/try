package try

// Scope executes f and converts any [Must]-induced panic back into an error
// returned via err.
//
// Panics whose value is not a panickedError are re-panicked unchanged, so real
// bugs (nil dereference, out-of-range, third-party panics) are never silently
// swallowed as errors.
//
// Only panics raised on the same goroutine that calls Scope are captured. If
// f spawns goroutines that themselves call Must, their panics will escape
// Scope and crash the program; a parent goroutine cannot reliably recover a
// child's panic.
//
// On the happy path (no panic) err is nil.
func Scope(f func()) (err error) {
	defer recoverScope(&err)
	f()
	return
}

// Scope1 is like [Scope] but f returns a single value, returned alongside the
// error. On a Must-induced panic, out is the zero value of A and err carries
// the original error. Foreign panics are re-panicked, as with Scope.
func Scope1[A any](f func() A) (out A, err error) {
	defer recoverScope(&err)
	out = f()
	return
}

// Scope2 is like [Scope1] but f returns two values. On a Must-induced panic,
// both outputs are zero-valued and err carries the original error.
func Scope2[A, B any](f func() (A, B)) (out1 A, out2 B, err error) {
	defer recoverScope(&err)
	out1, out2 = f()
	return
}

// Scope3 is like [Scope1] but f returns three values. On a Must-induced
// panic, all outputs are zero-valued and err carries the original error.
func Scope3[A, B, C any](f func() (A, B, C)) (out1 A, out2 B, out3 C, err error) {
	defer recoverScope(&err)
	out1, out2, out3 = f()
	return
}

// Scope4 is like [Scope1] but f returns four values. On a Must-induced panic,
// all outputs are zero-valued and err carries the original error.
func Scope4[A, B, C, D any](f func() (A, B, C, D)) (out1 A, out2 B, out3 C, out4 D, err error) {
	defer recoverScope(&err)
	out1, out2, out3, out4 = f()
	return
}

// Scope5 is like [Scope1] but f returns five values. On a Must-induced panic,
// all outputs are zero-valued and err carries the original error.
func Scope5[A, B, C, D, E any](f func() (A, B, C, D, E)) (out1 A, out2 B, out3 C, out4 D, out5 E, err error) {
	defer recoverScope(&err)
	out1, out2, out3, out4, out5 = f()
	return
}

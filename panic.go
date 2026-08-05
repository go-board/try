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

func recoverScope(err *error) {
	r := recover()
	if r == nil {
		return
	}
	pe, ok := r.(panickedError)
	if !ok {
		panic(r) // re-panic: not from our Must
	}
	*err = pe.err
}

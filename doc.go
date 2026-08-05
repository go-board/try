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
// Use [Of] / [Of1] .. [Of5] when a fallible call needs local context before it
// bubbles out. The returned Try value can wrap or map the error, then finish
// with Must inside a Scope or with Result at a normal Go return boundary:
//
//	func load(path string) (Config, error) {
//	    return try.Of1(readConfig(path)).Wrapf("read config %q", path).Result()
//	}
//
// Use [Assert] or [Assertf] for preconditions that should become ordinary
// errors at a Scope boundary.
//
// Only panics raised by this package's Must family are captured by Scope. Any
// other panic (nil dereference, out-of-range, third-party library panic) is
// re-panicked unchanged so genuine bugs are never silently swallowed.
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

package try

import (
	"errors"
	"fmt"
)

// ErrAssertionFailed is used when [Assert] fails without a specific error.
var ErrAssertionFailed = errors.New("assertion failed")

// Assert panics with err if ok is false.
//
// It is intended for use inside a [Scope] block, like [Must]. If ok is false
// and err is nil, Assert panics with [ErrAssertionFailed] so a failed assertion
// is not silently ignored.
func Assert(ok bool, err error) {
	if ok {
		return
	}
	if err == nil {
		err = ErrAssertionFailed
	}
	checkError(err)
}

// Assertf panics with a formatted error if ok is false.
//
// It is intended for use inside a [Scope] block, like [Assert].
func Assertf(ok bool, format string, args ...any) {
	if !ok {
		checkError(fmt.Errorf(format, args...))
	}
}

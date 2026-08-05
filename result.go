package try

import "fmt"

func mapError(err error, f func(error) error) error {
	if err == nil {
		return nil
	}
	mapped := f(err)
	if mapped == nil {
		return err
	}
	return mapped
}

func wrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	if message == "" {
		return err
	}
	return fmt.Errorf("%s: %w", message, err)
}

func wrapErrorf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return wrapError(err, fmt.Sprintf(format, args...))
}

// Try holds an error so it can be mapped or wrapped before [Try.Must] or
// [Try.Result].
type Try struct {
	err error
}

// Of creates a [Try] from an error-only result.
func Of(err error) Try {
	return Try{err}
}

// MapErr transforms the stored error. A nil mapper result keeps the original
// error so an existing failure is not swallowed accidentally.
func (t Try) MapErr(f func(error) error) Try {
	t.err = mapError(t.err, f)
	return t
}

// Wrap adds message as error context while preserving the original error for
// errors.Is and errors.As.
func (t Try) Wrap(message string) Try {
	t.err = wrapError(t.err, message)
	return t
}

// Wrapf formats error context and wraps the original error.
func (t Try) Wrapf(format string, args ...any) Try {
	t.err = wrapErrorf(t.err, format, args...)
	return t
}

// Must panics if the stored error is non-nil.
func (t Try) Must() {
	checkError(t.err)
}

// Result returns the stored error.
func (t Try) Result() error {
	return t.err
}

// Try1 holds one value and an error so the error can be mapped or wrapped
// before [Try1.Must] or [Try1.Result].
type Try1[A any] struct {
	v1  A
	err error
}

// Of1 creates a [Try1] from a value-plus-error result.
func Of1[A any](v1 A, err error) Try1[A] {
	return Try1[A]{v1, err}
}

// MapErr transforms the stored error. A nil mapper result keeps the original
// error so an existing failure is not swallowed accidentally.
func (t Try1[A]) MapErr(f func(error) error) Try1[A] {
	t.err = mapError(t.err, f)
	return t
}

// Wrap adds message as error context while preserving the original error for
// errors.Is and errors.As.
func (t Try1[A]) Wrap(message string) Try1[A] {
	t.err = wrapError(t.err, message)
	return t
}

// Wrapf formats error context and wraps the original error.
func (t Try1[A]) Wrapf(format string, args ...any) Try1[A] {
	t.err = wrapErrorf(t.err, format, args...)
	return t
}

// Must returns the stored value if the stored error is nil; otherwise it
// panics.
func (t Try1[A]) Must() A {
	checkError(t.err)
	return t.v1
}

// Result returns the stored value and error.
func (t Try1[A]) Result() (A, error) {
	return t.v1, t.err
}

// Try2 holds two values and an error so the error can be mapped or wrapped
// before [Try2.Must] or [Try2.Result].
type Try2[A, B any] struct {
	v1  A
	v2  B
	err error
}

// Of2 creates a [Try2] from a two-value-plus-error result.
func Of2[A, B any](v1 A, v2 B, err error) Try2[A, B] {
	return Try2[A, B]{v1, v2, err}
}

// MapErr transforms the stored error. A nil mapper result keeps the original
// error so an existing failure is not swallowed accidentally.
func (t Try2[A, B]) MapErr(f func(error) error) Try2[A, B] {
	t.err = mapError(t.err, f)
	return t
}

// Wrap adds message as error context while preserving the original error for
// errors.Is and errors.As.
func (t Try2[A, B]) Wrap(message string) Try2[A, B] {
	t.err = wrapError(t.err, message)
	return t
}

// Wrapf formats error context and wraps the original error.
func (t Try2[A, B]) Wrapf(format string, args ...any) Try2[A, B] {
	t.err = wrapErrorf(t.err, format, args...)
	return t
}

// Must returns the stored values if the stored error is nil; otherwise it
// panics.
func (t Try2[A, B]) Must() (A, B) {
	checkError(t.err)
	return t.v1, t.v2
}

// Result returns the stored values and error.
func (t Try2[A, B]) Result() (A, B, error) {
	return t.v1, t.v2, t.err
}

// Try3 holds three values and an error so the error can be mapped or wrapped
// before [Try3.Must] or [Try3.Result].
type Try3[A, B, C any] struct {
	v1  A
	v2  B
	v3  C
	err error
}

// Of3 creates a [Try3] from a three-value-plus-error result.
func Of3[A, B, C any](v1 A, v2 B, v3 C, err error) Try3[A, B, C] {
	return Try3[A, B, C]{v1, v2, v3, err}
}

// MapErr transforms the stored error. A nil mapper result keeps the original
// error so an existing failure is not swallowed accidentally.
func (t Try3[A, B, C]) MapErr(f func(error) error) Try3[A, B, C] {
	t.err = mapError(t.err, f)
	return t
}

// Wrap adds message as error context while preserving the original error for
// errors.Is and errors.As.
func (t Try3[A, B, C]) Wrap(message string) Try3[A, B, C] {
	t.err = wrapError(t.err, message)
	return t
}

// Wrapf formats error context and wraps the original error.
func (t Try3[A, B, C]) Wrapf(format string, args ...any) Try3[A, B, C] {
	t.err = wrapErrorf(t.err, format, args...)
	return t
}

// Must returns the stored values if the stored error is nil; otherwise it
// panics.
func (t Try3[A, B, C]) Must() (A, B, C) {
	checkError(t.err)
	return t.v1, t.v2, t.v3
}

// Result returns the stored values and error.
func (t Try3[A, B, C]) Result() (A, B, C, error) {
	return t.v1, t.v2, t.v3, t.err
}

// Try4 holds four values and an error so the error can be mapped or wrapped
// before [Try4.Must] or [Try4.Result].
type Try4[A, B, C, D any] struct {
	v1  A
	v2  B
	v3  C
	v4  D
	err error
}

// Of4 creates a [Try4] from a four-value-plus-error result.
func Of4[A, B, C, D any](v1 A, v2 B, v3 C, v4 D, err error) Try4[A, B, C, D] {
	return Try4[A, B, C, D]{v1, v2, v3, v4, err}
}

// MapErr transforms the stored error. A nil mapper result keeps the original
// error so an existing failure is not swallowed accidentally.
func (t Try4[A, B, C, D]) MapErr(f func(error) error) Try4[A, B, C, D] {
	t.err = mapError(t.err, f)
	return t
}

// Wrap adds message as error context while preserving the original error for
// errors.Is and errors.As.
func (t Try4[A, B, C, D]) Wrap(message string) Try4[A, B, C, D] {
	t.err = wrapError(t.err, message)
	return t
}

// Wrapf formats error context and wraps the original error.
func (t Try4[A, B, C, D]) Wrapf(format string, args ...any) Try4[A, B, C, D] {
	t.err = wrapErrorf(t.err, format, args...)
	return t
}

// Must returns the stored values if the stored error is nil; otherwise it
// panics.
func (t Try4[A, B, C, D]) Must() (A, B, C, D) {
	checkError(t.err)
	return t.v1, t.v2, t.v3, t.v4
}

// Result returns the stored values and error.
func (t Try4[A, B, C, D]) Result() (A, B, C, D, error) {
	return t.v1, t.v2, t.v3, t.v4, t.err
}

// Try5 holds five values and an error so the error can be mapped or wrapped
// before [Try5.Must] or [Try5.Result].
type Try5[A, B, C, D, E any] struct {
	v1  A
	v2  B
	v3  C
	v4  D
	v5  E
	err error
}

// Of5 creates a [Try5] from a five-value-plus-error result.
func Of5[A, B, C, D, E any](v1 A, v2 B, v3 C, v4 D, v5 E, err error) Try5[A, B, C, D, E] {
	return Try5[A, B, C, D, E]{v1, v2, v3, v4, v5, err}
}

// MapErr transforms the stored error. A nil mapper result keeps the original
// error so an existing failure is not swallowed accidentally.
func (t Try5[A, B, C, D, E]) MapErr(f func(error) error) Try5[A, B, C, D, E] {
	t.err = mapError(t.err, f)
	return t
}

// Wrap adds message as error context while preserving the original error for
// errors.Is and errors.As.
func (t Try5[A, B, C, D, E]) Wrap(message string) Try5[A, B, C, D, E] {
	t.err = wrapError(t.err, message)
	return t
}

// Wrapf formats error context and wraps the original error.
func (t Try5[A, B, C, D, E]) Wrapf(format string, args ...any) Try5[A, B, C, D, E] {
	t.err = wrapErrorf(t.err, format, args...)
	return t
}

// Must returns the stored values if the stored error is nil; otherwise it
// panics.
func (t Try5[A, B, C, D, E]) Must() (A, B, C, D, E) {
	checkError(t.err)
	return t.v1, t.v2, t.v3, t.v4, t.v5
}

// Result returns the stored values and error.
func (t Try5[A, B, C, D, E]) Result() (A, B, C, D, E, error) {
	return t.v1, t.v2, t.v3, t.v4, t.v5, t.err
}

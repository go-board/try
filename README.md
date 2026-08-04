# try

[![Go](https://github.com/go-board/try/actions/workflows/go.yml/badge.svg)](https://github.com/go-board/try/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/go-board/try/branch/main/graph/badge.svg)](https://codecov.io/gh/go-board/try)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-board/try.svg)](https://pkg.go.dev/github.com/go-board/try)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-board/try)](https://goreportcard.com/report/github.com/go-board/try)

`try` is a small Go package that provides a flat, exception-style error-handling
pattern on top of Go's native `error` values.

Instead of repeating `if err != nil { return err }` after every fallible call,
you write linear code using `Must*` inside a `Scope*` block. Errors are still
values at the block boundary — `Scope` recovers `Must`-induced panics and
returns them as ordinary errors.

## Why

Go's explicit error handling is great for readability at function boundaries,
but it pushes a lot of ceremony into the happy path:

```go
a, err := step1(in)
if err != nil {
    return 0, err
}
b, err := step2(a)
if err != nil {
    return 0, err
}
c, err := step3(b)
if err != nil {
    return 0, err
}
return c, nil
```

With `try` the same flow becomes linear, while the function still returns an
`error` rather than panicking:

```go
return try.Scope1(func() int {
    a := try.Must1(step1(in))
    b := try.Must1(step2(a))
    return try.Must1(step3(b))
})
```

## Installation

```
go get github.com/go-board/try
```

## Design

- **Errors stay values at the boundary.** `Scope`/`Scope1`..`Scope5` return a
  normal `error` (alongside any returned values). No `recover` leaks into
  caller code.
- **Only `Must` panics are captured.** A `Must` failure panics with an internal
  `panickedError` marker. `Scope` recovers only that marker; any other panic
  (nil dereference, out-of-range, third-party library panic) is re-panicked so
  real bugs are never silently swallowed.
- **Error identity is preserved.** The wrapped error round-trips through
  `errors.Is` / `errors.As`, including concrete types — see the example below.
- **Same-goroutine only.** Panics from goroutines spawned inside `f` are not
  captured — they will crash the program, which is the only safe default.
- **Zero values on failure.** When a `Must` panic is recovered by `Scope1`..
  `Scope5`, the returned values are the zero values of their types; only `err`
  is meaningful. Always check `err` before using the outputs.

### Error identity

`panickedError` implements `Unwrap`, so the original error survives the
panic → recover round-trip:

```go
var ErrEmpty = errors.New("empty input")

_, err := try.Scope1(func() int {
    return try.Must1(parse("")) // returns ErrEmpty
})

errors.Is(err, ErrEmpty)               // true
var ce *customErr
errors.As(err, &ce)                    // true if ErrEmpty wraps a *customErr
```

### Nested Scope

`Scope` blocks may be nested. The innermost `Scope` recovers a `Must` panic
first, so an inner failure does not propagate to an outer `Scope` unless the
inner `Scope` itself re-panics (which it never does for `Must` panics).

```go
err := try.Scope(func() {
    try.Must(outerStep()) // if this fails, outer Scope catches it

    if err := try.Scope(innerStep); err != nil {
        // inner failure handled here; outer Scope is unaffected
    }
})
```

## Pitfalls

- **Don't call `Must` in goroutines started inside `f`.** Their panics escape
  `Scope` and crash the program. If you need concurrent fallible work, run each
  in its own `Scope` and collect errors explicitly.
- **`Must` outside `Scope` crashes the program.** This is intentional (like
  `template.Must`); only use it for truly unrecoverable invariants.
- **Don't inspect outputs when `err != nil`.** They are zero values, not
  partial results.
- **Foreign panics are not converted.** A nil dereference or a panic from a
  third-party library inside `f` is re-panicked, not returned as `err`.

## API

### `Must` — assert no error

`Must` panics when `err != nil`. `Must1`..`Must5` additionally carry through
1–5 return values on the happy path.

```go
func Must(err error)
func Must1[A any](v A, err error) A
func Must2[A, B any](v1 A, v2 B, err error) (A, B)
func Must3[A, B, C any](v1 A, v2 B, v3 C, err error) (A, B, C)
func Must4[A, B, C, D any](v1 A, v2 B, v3 C, v4 D, err error) (A, B, C, D)
func Must5[A, B, C, D, E any](v1 A, v2 B, v3 C, v4 D, v5 E, err error) (A, B, C, D, E)
```

Outside a `Scope`, a `Must` panic propagates and crashes the program, similar
to `template.Must`.

### `Scope` — convert `Must` panics back to errors

`Scope` runs `f` and returns any `Must`-induced panic as an `err`.
`Scope1`..`Scope5` do the same when `f` returns 1–5 values.

```go
func Scope(f func()) (err error)
func Scope1[A any](f func() A) (out A, err error)
func Scope2[A, B any](f func() (A, B)) (out1 A, out2 B, err error)
func Scope3[A, B, C any](f func() (A, B, C)) (out1 A, out2 B, out3 C, err error)
func Scope4[A, B, C, D any](f func() (A, B, C, D)) (out1 A, out2 B, out3 C, out4 D, err error)
func Scope5[A, B, C, D, E any](f func() (A, B, C, D, E)) (out1 A, out2 B, out3 C, out4 D, out5 E, err error)
```

If a non-`Must` panic occurs inside `f`, `Scope` re-panics with the original
value so genuine bugs surface normally.

## Usage

```go
package main

import (
    "errors"
    "fmt"

    "github.com/go-board/try"
)

var ErrEmpty = errors.New("empty input")

func parse(s string) (int, error) {
    if s == "" {
        return 0, ErrEmpty
    }
    return len(s), nil
}

func double(s string) (int, error) {
    n := try.Must1(parse(s)) // panics on ErrEmpty
    return n * 2, nil
}

func main() {
    out, err := try.Scope1(func() int {
        return try.Must1(double("hello"))
    })
    fmt.Println(out, err) // 10 <nil>

    _, err = try.Scope1(func() int {
        return try.Must1(double(""))
    })
    fmt.Println(errors.Is(err, ErrEmpty)) // true
}
```

## When to use

- **Good fit:** short, linear setup or transformation pipelines where the only
  recovery strategy is "bubble the error up unchanged".
- **Not a good fit:** request handlers that need to map errors to status codes,
  retry loops, or any path where you actually branch on specific errors —
  plain `if err != nil` is clearer there.

## Compatibility

Requires Go 1.18+ (generics).

## License

See [LICENSE](LICENSE).

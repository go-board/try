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

When a fallible call needs local context before it bubbles out, wrap the result
with `Of*`, call `Wrap`, `Wrapf`, or `MapErr`, then finish with either `Must`
inside a `Scope` or `Result` at a normal Go return boundary.

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

`Of*` wrappers preserve the same identity when adding context:

```go
_, err := try.Scope1(func() int {
    return try.Of1(parse(input)).Wrap("parse input").Must()
})

errors.Is(err, ErrEmpty)               // true
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
  third-party library inside `f` is re-panicked with the original panic value,
  not returned as `err`.

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

### `Assert` — fail a precondition

`Assert` and `Assertf` are for preconditions that should bubble out as errors
from a `Scope`.

```go
var ErrEmpty = errors.New("empty input")

err := try.Scope(func() {
    try.Assert(input != "", ErrEmpty)
    try.Assertf(limit > 0, "invalid limit %d", limit)
})
```

```go
func Assert(ok bool, err error)
func Assertf(ok bool, format string, args ...any)
var ErrAssertionFailed error
```

If `Assert` fails with a nil error, it uses `ErrAssertionFailed` so the failed
assertion is not silently ignored. Outside a `Scope`, failed assertions panic
like `Must`.

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

Go does not have Python-style bare re-raise. The re-panic keeps the original
panic value, but runtime output or outer recover middleware may show the
`Scope` re-panic frame above the original business frames.

### `Of` — add context before returning

`Of`/`Of1`..`Of5` capture a normal Go result so the error can be wrapped or
mapped before the final `Must` or `Result`.

```go
func Of(err error) Try
func Of1[A any](v1 A, err error) Try1[A]
func Of2[A, B any](v1 A, v2 B, err error) Try2[A, B]
func Of3[A, B, C any](v1 A, v2 B, v3 C, err error) Try3[A, B, C]
func Of4[A, B, C, D any](v1 A, v2 B, v3 C, v4 D, err error) Try4[A, B, C, D]
func Of5[A, B, C, D, E any](v1 A, v2 B, v3 C, v4 D, v5 E, err error) Try5[A, B, C, D, E]
```

Each `Try*` value supports the same chainable methods:

```go
MapErr(func(error) error) Try*
Wrap(message string) Try*
Wrapf(format string, args ...any) Try*
Must() values...
Result() values..., error
```

`Wrap` and `Wrapf` add context while preserving the original error for
`errors.Is` and `errors.As`. `MapErr` is for custom error types or policies.
If `MapErr` returns nil for a non-nil error, the original error is kept so an
existing failure is not swallowed accidentally.

Use `Must` inside `Scope` when you want linear control flow. Use `Result` when
you are already at a normal Go return boundary.

```go
return try.Scope1(func() int {
    cfg := try.Of1(readConfig(path)).Wrapf("read config %q", path).Must()
    return try.Of1(parseConfig(cfg)).Wrap("parse config").Must()
})
```

```go
func load(path string) (Config, error) {
    return try.Of1(readConfig(path)).Wrapf("read config %q", path).Result()
}
```

### `tryerr` — structured error metadata

The `tryerr` subpackage adds stable codes, structured attributes, and a
backtrace while keeping Go's native error chain intact.

```go
return try.Of1(readConfig(path)).
    MapErr(func(err error) error {
        return tryerr.Wrap(
            err,
            "read config",
            tryerr.WithCode(1001),
            tryerr.WithAttr("path", path),
        )
    }).
    Result()
```

`tryerr.Wrap(nil, ...)` returns nil, so it is safe to use in ordinary
propagation paths. The default code is `-1`; override it with
`tryerr.WithCode`. Code `0` is reserved for unset values and is skipped by
`tryerr.Code`. `tryerr.WithAttrs` copies its input map, and `tryerr.Attrs`
returns a fresh merged map where outer errors keep precedence on key collisions.

```go
func New(message string, opts ...Option) error
func Wrap(cause error, message string, opts ...Option) error
func WithCode(code int) Option
func WithAttr(key, value string) Option
func WithAttrs(attrs map[string]string) Option
func WithStackDepth(depth int) Option
func WithStackSkip(skip int) Option
func Code(err error) int
func Attrs(err error) map[string]string
func StackFrames(err error) iter.Seq[runtime.Frame]
```

Use `errors.Is` / `errors.As` for the original cause and `tryerr.Code`,
`tryerr.Attrs`, or `tryerr.StackFrames` for structured metadata.
`StackFrames` returns `iter.Seq[runtime.Frame]`, so callers can range over
resolved frames without handling raw program counters. `New` and `Wrap` capture
one caller frame by default; use `tryerr.WithStackDepth` to capture more and
`tryerr.WithStackSkip` to skip additional caller frames. Stack skip is relative
to the caller of `New` or `Wrap`. `Error()` includes the code and first captured
frame in the returned text, for example `read config [code=1001] [...]`.

The extractor functions primarily read `tryerr.Error` values. They also accept
small same-shape providers in the chain: `interface{ Code() int }`,
`interface{ Attrs() map[string]string }`, and
`interface{ StackFrames() iter.Seq[runtime.Frame] }`.

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
    return try.Scope1(func() int {
        n := try.Must1(parse(s)) // converted to err by Scope1
        return n * 2
    })
}

func main() {
    out, err := double("hello")
    fmt.Println(out, err) // 10 <nil>

    _, err = double("")
    fmt.Println(errors.Is(err, ErrEmpty)) // true
}
```

## When to use

- **Good fit:** short, linear setup or transformation pipelines where the only
  recovery strategy is "bubble the error up", optionally with local context.
- **Not a good fit:** request handlers that need to map errors to status codes,
  retry loops, or any path where you actually branch on specific errors —
  plain `if err != nil` is clearer there.

## Compatibility

Requires Go 1.24+.

## License

See [LICENSE](LICENSE).

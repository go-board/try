# Contributing Guide

Thanks for contributing to `go-board/try`.

## Development style

### Language and formatting

- Use English for all public comments, docs, and test descriptions.
- Keep comments concise and behavior-focused.
- Run `gofmt` on all changed Go files before submitting.

### API conventions

- Keep the two families aligned by arity: `Must0`..`Must5` pair with
  `Scope0`..`Scope5`.
- `Must*` signatures take `(values..., err)` with `err` last. Do not break
  this ordering.
- All `Must*` variants delegate to the single `checkError` chokepoint so the
  panic origin stays consistent.
- `Scope*` recovers only `panickedError`; any other panic is re-panicked so
  real bugs surface.
- Error identity must round-trip through `errors.Is` / `errors.As` via
  `panickedError.Unwrap` — do not flatten the wrapped error.

### Tests

- Add/adjust tests for any public API change.
- Prefer table-driven tests when multiple edge-cases share setup.
- Keep assertion messages explicit and action-oriented.
- Cover both the happy path (err == nil) and the failure path (panic captured
  by `Scope`), plus the re-panic behavior for foreign panics.

## Validation checklist

- `gofmt -w <changed_files>`
- `go vet ./...`
- `go test -race ./...`

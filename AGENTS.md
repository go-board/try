# AGENTS.md

This file defines repository-wide instructions for coding agents working in this project.

## Scope

These instructions apply to the entire repository.

## Language and comments

- Use English for all public-facing comments, docs, and test descriptions.
- Keep comments concise and behavior-focused; explain "why", not "what".
- Remove stale comments during refactors.

## API style

- The package exposes two families: `Must`/`Must1`..`Must5` and `Scope`/`Scope1`..`Scope5`.
  Keep numbering aligned across both families (e.g. `Must3` pairs with `Scope3`).
- `Must*` takes `(values..., err)` with `err` always last. Do not introduce
  positional variants that break this convention.
- All `Must*` variants must delegate to the single `checkError` chokepoint so
  the panic origin and `panickedError` wrapping stay consistent.
- `Scope*` must recover only `panickedError`; any other panic is re-panicked
  unchanged so real bugs are never silently swallowed.
- Preserve error identity through `panickedError.Unwrap` — do not flatten or
  rewrap the original error.

## Error handling

- Never add a global error callback or global formatter that mutates the Must
  hot path. If context wrapping is needed, do it at the `Scope` boundary in
  caller code (`fmt.Errorf("...: %w", err)`).
- Document that `Must` outside `Scope` crashes the program (like `template.Must`).

## Testing and validation

- Add or update tests for behavior changes.
- Run the following before finalizing:
  - `gofmt -w <changed_files>`
  - `go vet ./...`
  - `go test -race ./...`

## Documentation

- Keep `README.md` aligned with public API changes.
- Prefer examples that compile against current APIs.

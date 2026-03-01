# Contributing to Kwelea

## Prerequisites

- Go (see `go.mod` for the minimum version)
- [just](https://just.systems) -- task runner
- [lefthook](https://github.com/evilmartians/lefthook) -- git hooks
- [gotestsum](https://github.com/gotestyourself/gotestsum) -- test runner

```bash
go install gotest.tools/gotestsum@latest
lefthook install
```

## Development workflow

```bash
just build        # compile the binary to bin/
just test         # run tests
just check        # format + vet
just serve        # build and serve kwelea's own docs (dogfood)
```

Pre-commit hooks run `fmt`, `vet`, and `test` on changed `.go` files.

## Making changes

For non-trivial changes, open an [issue](https://github.com/engineervix/kwelea/issues) first to discuss the approach. Bug fixes and typos can go straight to a PR.

If you are proposing a feature, explain what problem it solves and keep the scope narrow — a focused PR is easier to review and more likely to land.

1. Fork the repo and create a branch (`feat/my-thing`, `fix/my-bug`).
2. Write tests for new behaviour. Existing tests live in `internal/*/` alongside the code.
3. Update `docs/` alongside the code if the change affects users.
4. Open a pull request against `main`.

## Commit format

This project uses [Conventional Commits](https://www.conventionalcommits.org). The release changelog is generated automatically from commit messages, so the format matters.

```
feat: add search highlighting for multi-word queries
fix: resolve port conflict on dev server restart
docs: clarify nav ordering in configuration guide
chore: bump D2 to v0.7.2
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`.
Breaking changes: append `!` after the type (`feat!:`) and describe the break in the commit body.

## Running the full test suite

```bash
just test-v          # verbose output
```

CI runs on Linux and macOS. Windows builds are cross-compiled in CI -- flag any platform-specific issues in your PR.

# AGENTS.md

Instructions for agentic coding tools working in this repository.

## Project Overview

- Language: Go (`go 1.26.1` in `go.mod`).
- App type: Bubble Tea terminal UI for managing Zellij sessions.
- Entrypoint: `cmd/zema/main.go`.
- Main packages:
  - `internal/view`: TUI model (`Init/Update/View`), forms, key handling, styling.
  - `internal/zema`: app composition and startup.
  - `pkg/zellij`: wrapper around the `zellij` CLI.
  - `pkg/cmd`: timeout-aware command execution helper.

## Repository Notes

- Build system: `Taskfile.yml` (primary task: `build`).
- Build output: `bin/zema`.
- Dependency source snapshots: `opensrc/`.
- Cursor rules:
  - `.cursor/rules/` not present.
  - `.cursorrules` not present.
- Copilot rules:
  - `.github/copilot-instructions.md` not present.

## Build / Run / Lint / Test

### Run locally

```bash
go run .
```

### Build

Preferred reproducible build:

```bash
task build
```

Simple module build:

```bash
go build ./...
```

### Lint and static checks

No repository linter config is defined (no root `.golangci.yml`).
Use Go-native checks:

```bash
gofmt -w .
go vet ./...
```

Format-check only:

```bash
gofmt -l .
```

## Coding Conventions

### Imports

- Use standard Go grouping:
  1. standard library
  2. blank line
  3. third-party and module-local imports
- Use module-local paths like `zema/internal/view`, `zema/pkg/cmd`.
- Keep imports minimal; remove unused imports in every change.

### Formatting

- Always run `gofmt` after editing Go code.
- Keep functions small and purpose-specific.
- Avoid unnecessary comments; rely on clear names and structure.
- Reuse existing constants instead of duplicating literals.

### Types and package design

- Prefer concrete structs for stateful logic (`model`, `Zema`, `Cmd`).
- Introduce interfaces only at real boundaries (example: `view.Zellij`).
- Keep internals unexported unless cross-package use is required.
- Return errors instead of panicking for expected runtime failures.

### Naming

- Exported identifiers: `CamelCase`.
- Unexported identifiers: `camelCase`.
- Short receiver names are preferred (`m`, `z`, `cmd`).
- Group related constants in `const` blocks (`internal/view/const.go`).
- Existing code contains some underscore-prefixed locals (for example `_zellij`).
  - Preserve local consistency when editing an existing file.
  - In new code, prefer idiomatic non-underscore names unless needed.

### Error handling

- Wrap errors with context and `%w`.
  - Example: `fmt.Errorf("failed to create session: %w", err)`.
- Keep user-facing messages actionable and concise.
- For external command failures, include command output context where useful.
- Preserve timeout behavior in `pkg/cmd` (`ErrTimeout`).

### Bubble Tea / TUI patterns

- Keep `Init`, `Update`, and `View` responsibilities clearly separated.
- In `Update`, handle message type first, then state/view branching.
- Prefer explicit state transitions (`m.view = viewAdd`).
- Rebuild forms when underlying data changes (`refreshSessionsForm`, `refreshAddForm`).
- Keep key bindings centralized in `internal/view/const.go`.

### Command execution

- Use `pkg/cmd` for external command execution with timeout/cancel support.
- Always `defer cancel()` for timeout contexts.
- Capture combined output for diagnostics (`CombinedOutput`).

## Test Guidance For New Code

- Add table-driven tests for branch-heavy logic.
- Put tests near the package being changed.
- For `pkg/cmd`, cover timeout and command failure paths.
- For `internal/view`, test key flows and resulting model state changes.

## Expected Agent Workflow

- Before editing: inspect nearby files and keep conventions consistent.
- After editing: run `go test ./...` and `go vet ./...`.
- If UI behavior changes, smoke test with `go run .`.
- Keep diffs focused; avoid unrelated refactors.

<!-- opensrc:start -->

## Source Code Reference

Source code for dependencies is available in `opensrc/` for deeper understanding of implementation details.

See `opensrc/sources.json` for the list of available packages and their versions.

Use this source code when you need to understand how a package works internally, not just its types/interface.

### Fetching Additional Source Code

To fetch source code for a package or repository you need to understand, run:

```bash
opensrc <package>           # npm package (e.g., opensrc zod)
opensrc pypi:<package>      # Python package (e.g., opensrc pypi:requests)
opensrc crates:<package>    # Rust crate (e.g., opensrc crates:serde)
opensrc <owner>/<repo>      # GitHub repo (e.g., opensrc vercel/ai)
```

<!-- opensrc:end -->

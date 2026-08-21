# cmd/solast — CLI

## Purpose

Cobra CLI wrapper over [[parser-index]]. Useful for manual inspection and as a reference for the public API. Built via `make build` (embeds version/commit/build-time through ldflags).

## main.go

**Build vars** (main.go:15): `Version`, `BuildTime`, `GitCommit` — set by ldflags, else from module build info.

**Subcommands:**
- `parse [file|-]` (main.go:63) → JSON AST. Flags: `--output/-o`, `--loc`, `--range`, `--tolerant`, `--pretty/-p` (default true). Handler `runParse` (106).
- `validate [file|-]` (main.go:79) → syntax check; exit 0 valid / 1 on errors; errors to stderr as `line:column: message`. Handler `runValidate` (136), tolerant internally.
- `version-detect [file|-]` (main.go:89) → prints detected pragma/version/constraint. Handler `runVersionDetect` (162).

**Helpers:** `readInput` (182, file or stdin), `writeOutput` (204, file or stdout + trailing newline).

**Root** (main.go:53): `Use: "solast"`, version string `X.Y.Z (commit: …, built: …)`.

## When this changes

Add a subcommand here only when exposing a new parser capability on the CLI. Most syntax work needs NO change here — `parse`/`validate` exercise the full parser generically. Keep flag names stable.

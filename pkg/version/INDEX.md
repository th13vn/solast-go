# pkg/version — Solidity Version Detection

## Purpose

Parse and compare **Solidity language versions** and detect `pragma solidity` constraints from source. Independent of the module's own release version (that's the `VERSION` file). Used to reason about which language features a contract targets.

## version.go

**Version** (version.go:12): `{Major, Minor, Patch int}`.
- `New(major, minor, patch) Version` (19), `String()` (24).
- Comparisons: `Compare(other) int` (30, -1/0/1), `LessThan`/`LessThanOrEqual`/`GreaterThan`/`GreaterThanOrEqual`/`Equal` (53-75), `IsZero()` (78).

**Parsing:**
- `Parse(s string) (Version, error)` (83) — `"0.8.20"` or `"0.8"` (patch defaults to 0).
- `MustParse(s)` (112) — panics on error.
- `ParseConstraint(s) (constraint string, v Version, err error)` (123) — splits a leading `^ ~ >= <= > < =` operator from the version. Regex `^(\^|~|>=|<=|>|<|=)?(\d+\.\d+(\.\d+)?)$`.

**Detection from source:**
- `Detect(source) (*DetectedVersion, error)` (148) — first `pragma solidity …;`.
- `DetectAll(source) ([]*DetectedVersion, error)` (184) — all pragmas; skips malformed.
- `DetectedVersion{Raw, Constraint string; Version Version}` (140).

## When this changes

Only when Solidity's version/pragma syntax changes (rare). Adding a new minor version of the language does NOT require changes here — it parses generically. New constraint operators would extend `ParseConstraint`'s regex.

## Tests
`version_test.go` — parse/compare/detect tables.

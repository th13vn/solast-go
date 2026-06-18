# scripts

## generate.sh

Generates a Go ANTLR parser from `grammar/*.g4` into `internal/gen/` (downloads `antlr-4.13.1-complete.jar`, requires Java). **Reference/experimental only** — the generated code is NOT used by the build; the runtime parser is hand-written in [[builder]]. Invoked by `make generate`.

See [[grammar-index]] for the grammar↔parser relationship and `make update-grammar` (which pulls the upstream `.g4`).

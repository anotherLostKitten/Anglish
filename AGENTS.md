# Repository Guidelines

## Project Structure & Module Organization

CLI entry lives in `cmd/anglish/main.go`, which wires the parser and compiler
together. Core packages sit under `internal/`: `parse` handles `.ang` syntax,
`compile` produces runnable plans, and `llm` wraps LangChainGo integrations.
Reference prompts in `examples/`, long-form notes in `docs/`, and shared test
data in `tests/`. Environment defaults reside in `env.example`; copy it to
`.env` when running locally.

## Build, Test, and Development Commands

- `go run ./cmd/anglish/main.go examples/chatbot.ang` runs the CLI against the
  sample program (pass your own `.ang` path to iterate quickly).
- `go build ./cmd/anglish` produces the standalone `anglish` binary for
  integration testing.
- `go test ./...` executes fast parser/compiler unit tests; add `-v` when
  debugging failures.
- `go test -tags vllm ./tests` enables the LLM integration suite that talks to
  the endpoint configured via `OPENAI_BASE_URL`.
- `go vet ./...` catches common Go pitfalls; run it before every PR to keep
  reviews short.

## Coding Style & Naming Conventions

Format code with `gofmt -w` (tabs for indentation, goimports-style imports).
Stick to idiomatic Go naming: exported APIs use PascalCase, private helpers use c
amelCase, and test files end with `_test.go`. Keep packages cohesive—avoid
cross-importing between `internal/parse` and `internal/compile` unless the
dependency direction is explicit in the parser output structs.

## Testing Guidelines

Unit tests rely on `stretchr/testify/require` for assertions; follow that style
for consistency. Name tests `Test<ThingUnderTest>` and co-locate them beside the
package they cover unless they need the shared harness in `tests/`. Integration
tests behind the `vllm` build tag require a reachable OpenAI-compatible server;
skip or mark flaky cases instead of removing them. Aim to cover new parsing
branches and compiler passes, and add fixtures under `examples/` when it
clarifies behavior.

## Commit & Pull Request Guidelines

Existing history favors short, imperative commit subjects
("flip dependency direction"), so do the same and explain the why in the body if
needed. Every PR should reference the affected module(s), describe test coverage
(`go test`, `go test -tags vllm`, etc.), and link to any tracking issue. Include
screenshots or transcripts when UI- or agent-facing behavior changes, and call
out new env vars so reviewers can update their `.env` files.

## Environment & Security Tips

Load configuration via `godotenv`; never commit secrets. Use `.env` locally, but
rely on real env vars in CI. Keep `OPENAI_API_KEY`, `HUGGING_FACE_HUB_TOKEN`, and
`OPENAI_MODEL` in sync with the backend you test against, and prefer mockable
interfaces when touching `internal/llm` to avoid leaking credentials into logs.

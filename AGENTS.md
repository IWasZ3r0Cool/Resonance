# Resonance contributor instructions

## Mission

Resonance is a local-first resume engine. Preserve candidate truth, privacy, and reproducibility before optimizing generated prose.

## Repository map

- `cmd/api`: process startup and graceful shutdown only.
- `internal/config`: environment parsing and validation.
- `internal/database`: SQLite initialization and forward-only migrations.
- `internal/agentengine`: Google ADK construction and model-provider boundaries.
- `internal/httpapi`: REST handlers and HTTP middleware.
- `api/openapi.json`: source of truth for the public API.
- `web`: React 19 + TypeScript + Vite client.
- `.agents/skills`: reusable project-specific workflows.

## Non-negotiable invariants

1. Never invent or silently strengthen a candidate fact. Generated claims must trace to candidate-provided evidence.
2. Preserve the exact input snapshot and generation metadata for every resume used in an application.
3. Keep model-vendor types inside `internal/agentengine`; domain and transport code depend on project-owned types.
4. Keep credentials server-side. No API key may use a `VITE_` prefix or appear in API responses/logs.
5. Bind to loopback by default. Treat LAN/public exposure as a separate security feature.
6. SQLite migrations are forward-only. Do not edit a migration that may already have shipped; add a new one.
7. `api/openapi.json` is contract authority. Update it, backend behavior, generated frontend types, and callers atomically.

## Relevant skills

- Invoke `resonance-domain` for candidate data, gap analysis, resume generation, scoring, or application tracking changes.
- Invoke `resonance-openapi` for any added or changed HTTP endpoint, schema, or frontend API call.

## Commands

```bash
make bootstrap
make generate-api
make test
make build
make check
```

Run focused tests while iterating, then run `make check` before handing off cross-stack changes. Tests must not require a real model key or network service unless they are explicitly marked as integration tests.

- Go unit tests live beside their package and run with `go test ./...`.
- Cross-package integration tests live in `tests/integration` and run with `go test -tags=integration ./tests/integration`.
- Frontend unit tests use Vitest and Testing Library and run with `pnpm --dir web test`.
- CI additionally runs the Go race detector, coverage gates, CodeQL, dependency review, and vulnerability scans.

## Go conventions

- Keep `cmd` thin and put behavior in `internal` packages.
- Pass `context.Context` across I/O boundaries.
- Wrap errors with the operation that failed; do not log and return the same error.
- Prefer standard-library facilities; add dependencies only when they remove meaningful maintenance or correctness risk.
- Do not initialize an LLM at process startup if that would make ordinary local development require credentials.

## Frontend conventions

- Use strict TypeScript and generated OpenAPI types. Do not duplicate server DTOs by hand.
- Use Material UI for interactive/accessibility-sensitive components and Tailwind for layout utilities.
- Keep data access under `web/src/api` rather than calling `fetch` throughout components.
- Represent loading, empty, ready, and error states explicitly.

## Definition of done

- Relevant tests and type checks pass.
- The OpenAPI document parses and generated types are current.
- No candidate truth, credential, or local-only security invariant was weakened.
- New configuration is documented in `.env.sample` and README without committing secrets.

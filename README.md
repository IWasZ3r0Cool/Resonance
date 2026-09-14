# Resonance Resume Engine

Resonance is a local-first, open-source workspace for turning a complete, truthful record of a candidate's work into a focused resume for each job application.

Instead of repeatedly rewriting a resume, the candidate maintains a private source of truth: roles, projects, skills, evidence, and accomplishments can be captured in as much detail as necessary. Resonance then selects and compresses the most relevant facts for a job description, identifies genuine gaps, and preserves the exact resume and job context used for every application.

> [!IMPORTANT]
> Resonance may rephrase, prioritize, and condense candidate-provided facts. It must never invent employment, credentials, skills, dates, metrics, or accomplishments.

## Project status

This repository is at the foundation stage. The current stack provides:

- a Go REST API with SQLite persistence and a checked-in OpenAPI 3.1 contract;
- a React 19 + TypeScript client built with Vite, Material UI, and Tailwind CSS;
- a typed frontend API client generated from `api/openapi.json`;
- Google ADK for Go wired behind a provider-neutral model factory;
- OpenAI, Gemini, and local/OpenAI-compatible model configuration;
- initial candidate, experience, skill, job, resume-version, and application tables;
- repository instructions and reusable Codex skills under `.agents/skills`.

The health and capabilities endpoints are runnable. Candidate editing, generation workflows, document rendering, and application management are the next product slices.

## Product principles

1. **The candidate record is the source of truth.** Generated content must be traceable to stored facts.
2. **Gaps are information, not invitations to fabricate.** Ask the candidate for missing evidence or produce the best honest draft.
3. **Every application is reproducible.** Draft resume versions are editable, but final versions are immutable snapshots of the job description, selected evidence, generation settings, and rendered resume.
4. **Local-first means useful without a cloud platform.** SQLite is the default and recent Ollama, LM Studio, and vLLM deployments can be used through an OpenAI Responses-compatible endpoint.
5. **Providers are replaceable.** Domain and workflow code must not depend directly on a model vendor.

## Architecture

```text
React 19 / TypeScript
  Material UI + Tailwind
          |
          | generated OpenAPI client
          v
Go REST API (chi) ------ api/openapi.json
    |          |
    |          +------ Google ADK for Go
    |                    | OpenAI
    |                    | Gemini
    |                    ` OpenAI-compatible local endpoint
    v
SQLite (local source of truth)
```

The browser talks only to the Resonance API. The backend owns persistence, truth-preserving workflow rules, provider credentials, and agent execution. API keys are never exposed through `VITE_*` variables or sent to the browser.

## Stack decisions

### Frontend: Vite + React + TypeScript

Use TypeScript. The OpenAPI contract can generate compile-time API types, and the extra safety is particularly valuable for structured candidate evidence and immutable resume snapshots.

Vite is the recommended bundler here. It is a better fit than a server-rendered React framework because Resonance is a local application with a separate Go API: Vite has a small runtime surface, fast development startup, efficient production output, and no duplicate Node server to operate. React 19, Material UI, and Tailwind CSS all have first-class Vite integrations.

Material UI provides accessible application primitives. Tailwind is reserved for layout and small visual adjustments; avoid building two competing theme systems.

### Agent SDK: Google ADK for Go

The recommendation is **Google ADK for Go v2**, currently pinned in `go.mod`. It is the only option among Google ADK, OpenAI Agents SDK, and NVIDIA NeMo Agent Toolkit that provides a maintained, native-Go orchestration SDK while also supporting the provider portability this project needs.

Top five advantages:

1. Native, idiomatic Go keeps the agent runtime in the backend process and avoids a Python sidecar.
2. It supplies agents, typed tools, sessions, workflow graphs, parallel/loop execution, and human-in-the-loop confirmation.
3. Its `model.LLM` boundary keeps orchestration separate from the model provider.
4. Its OpenAI Responses adapter works with OpenAI and recent compatible Ollama, LM Studio, and vLLM servers through a base URL.
5. It is open source, testable as ordinary Go code, and deployable anywhere a Go service runs.

Top five tradeoffs:

1. The Go implementation is younger than ADK Python and still changes more quickly, so upgrades need review.
2. Some integrations and examples remain Python-first.
3. The OpenAI adapter is marked experimental, even though it is the most direct path to local provider compatibility.
4. Responses-API compatibility is required; servers that only emulate Chat Completions cannot use that adapter.
5. The dependency graph is relatively large for a Go service because ADK includes broad orchestration and provider capabilities.

Why not the alternatives?

- **OpenAI Agents SDK:** excellent minimal primitives, guardrails, handoffs, tracing, and evaluation, but its mature SDKs are Python and TypeScript. OpenAI's official Go client exposes the Responses and beta Agents APIs, but it is not the same in-process Go orchestration SDK.
- **NVIDIA NeMo Agent Toolkit:** strong framework-agnostic observability and enterprise integration, but it is a Python library. Adopting it now would add a second service language without improving the local-first Go architecture.

Provider creation is isolated in `internal/agentengine`. If the ecosystem changes, Resonance can replace the SDK without contaminating candidate or application domain code.

## Quick start

### Prerequisites

- Go 1.26.6 or newer (required by the pinned ADK release)
- Node.js 22 or newer
- pnpm 11 (the repository pins the package-manager version)
- optionally, an OpenAI or Gemini key, or a local Responses-compatible model server

### Install and run

```bash
cp .env.sample .env
pnpm --dir web install --frozen-lockfile
go run ./cmd/api
```

In another terminal:

```bash
pnpm --dir web dev
```

Open <http://localhost:5173>. The Vite development server proxies `/api` and `/openapi.json` to the Go API at <http://localhost:8080>.

The API starts without a model key so the UI and database can be developed independently. Agent execution will remain unavailable until the selected provider is configured.

### Local model example

Recent versions of Ollama expose an OpenAI-compatible Responses API. After pulling a tool-capable model, update `.env`:

```dotenv
LLM_PROVIDER=ollama
LLM_MODEL=qwen3:8b
LLM_BASE_URL=http://localhost:11434/v1
LLM_API_KEY=ollama
```

The model must support the structured output and tool-calling behavior required by the workflow. A successful HTTP connection alone is not proof that a model is capable enough to tailor resumes safely.

## Configuration

Copy `.env.sample` to `.env`. The server loads `.env` for local development, while real environment variables take precedence.

| Variable | Default | Purpose |
| --- | --- | --- |
| `APP_ENV` | `development` | Runtime environment label. |
| `HTTP_HOST` | `127.0.0.1` | API bind address. Use `0.0.0.0` only when LAN/container access is intended. |
| `HTTP_PORT` | `8080` | API port. |
| `WEB_ORIGIN` | `http://localhost:5173` | Allowed browser origin. |
| `DATABASE_URL` | local `data/resonance.db` | SQLite DSN. |
| `LLM_PROVIDER` | `openai` | `openai`, `gemini`, `ollama`, or `openai-compatible`. |
| `LLM_MODEL` | `gpt-5-mini` | Provider model identifier. |
| `LLM_BASE_URL` | empty | Required for a custom OpenAI-compatible endpoint; optional preset for Ollama. |
| `LLM_API_KEY` | empty | Key for a generic compatible endpoint. |
| `OPENAI_API_KEY` | empty | OpenAI credential, read only by the backend. |
| `GOOGLE_API_KEY` | empty | Gemini credential, read only by the backend. |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, or `error`. |
| `MAX_UPLOAD_BYTES` | `10485760` | Future resume/import upload limit. |
| `VITE_API_BASE_URL` | empty | Browser API prefix; empty uses the same origin/proxy. Never put secrets here. |

## API contract

`api/openapi.json` is the source of truth for the HTTP interface.

```bash
pnpm --dir web generate:api  # refresh web/src/api/schema.d.ts
pnpm --dir web typecheck
```

Any API change must update the contract, backend behavior/tests, generated frontend types, and frontend call sites in the same change. The initial endpoints are:

- `GET /api/v1/health`
- `GET /api/v1/meta/capabilities`
- `GET /openapi.json`

## Development commands

```bash
make bootstrap      # install frontend dependencies and generate API types
make dev-api        # run the Go API
make dev-web        # run Vite
make generate-api   # regenerate TypeScript types
make test           # Go tests and frontend checks
make build          # compile API and production web assets
make check          # formatting, tests, and builds
```

### Verification layers

- Go package-level unit tests cover configuration, provider construction, migrations, contract embedding, and handlers.
- Vitest + Testing Library cover browser states and the typed API client with enforced coverage thresholds.
- `go test -tags=integration ./tests/integration` brings up a freshly migrated SQLite database behind the real router and verifies the cross-package API behavior without model credentials.
- GitHub Actions runs formatting, lint, type checks, `go vet`, the Go race detector, both builds, generated-contract drift detection, CodeQL, dependency review, `govulncheck`, and a high-severity npm audit.
- Dependabot opens grouped Go/frontend updates and monthly workflow-action updates.

Repository maintainers should also enable GitHub secret scanning, push protection, private vulnerability reporting, and branch protection requiring the CI, Security, and CodeQL checks. Those controls are repository settings and cannot be guaranteed by committed workflow files alone.

## Dependency policy

Direct runtime and build dependencies are pinned in `go.mod`, `web/package.json`, and `web/pnpm-lock.yaml`; transitive Go dependencies are recorded in `go.sum`.

| Area | Dependency | Role |
| --- | --- | --- |
| Backend | Go | API, workflows, and local runtime |
| Routing | `go-chi/chi` + `go-chi/cors` | Small HTTP router and explicit browser-origin policy |
| Agent runtime | `google.golang.org/adk/v2` | Provider-neutral agent and workflow orchestration |
| Database | `modernc.org/sqlite` | CGO-free local SQLite driver |
| Local env | `joho/godotenv` | Loads `.env` without overriding process variables |
| UI | React 19 + Material UI | Component application shell |
| Styling | Tailwind CSS | Layout and small utility-level styling |
| Build | Vite + TypeScript | Fast browser build and strict type checking |
| API client | `openapi-typescript` + `openapi-fetch` | Contract-generated types and typed requests |

Review dependency updates deliberately. ADK, the generated client, and database drivers should never be floated to unreviewed versions in production images.

TypeScript is intentionally pinned to the latest compatible 5.x release: the current `openapi-typescript` generator requires TypeScript 5, so upgrading the compiler to 6/7 must wait for that toolchain to declare compatibility.

## Initial data model

- `candidate_profiles`: the owner of a local workspace.
- `experiences`: detailed positions and projects, including source notes.
- `resume_bullets`: candidate-owned, reusable accomplishment evidence with optional situation/action/result detail; a rendered résumé may condense it without mutating it.
- `skills` + `bullet_skills`: normalized skills such as Go or React, linked to the exact same-candidate bullets that demonstrate them with provenance.
- `tags` + `bullet_tags`: typed hats and differentiators—such as leadership, technical excellence, domain, or outcome—linked within the same candidate profile with provenance and rationale.
- `job_targets`: immutable job-description snapshots and source metadata.
- `resume_versions`: editable drafts and immutable final artifacts containing the factual source snapshot and generation metadata. A revision of a final artifact is a new draft whose `parent_version_id` points to that finalized same-candidate version.
- `applications`: the company/role/status record linked only to the exact finalized resume version used.

Skills and tags are many-to-many: one bullet can demonstrate several attributes, and one attribute can be supported by several bullets. Composite foreign keys prevent associations from crossing candidate profiles. Association `source` values distinguish candidate-provided/imported metadata from model suggestions that the candidate has explicitly confirmed. Unconfirmed inference must remain outside the source-of-truth tables; a model may identify a gap, but it may not invent or silently strengthen evidence to close one.

SQLite foreign keys, WAL mode, and a busy timeout are enabled by the default DSN. Migrations run at API startup and are designed to move forward only.

The foundation migration is not yet released and may still change before the first stable schema. If a development database was created from an earlier copy of `001_initial.sql`, delete and recreate that local database; there is no in-place upgrade path for prerelease foundation data.

## Repository guide for coding agents

See `AGENTS.md`. Codex-compatible project skills live in:

- `.agents/skills/resonance-domain/SKILL.md`
- `.agents/skills/resonance-openapi/SKILL.md`

These skills capture the two non-obvious constraints that must survive future implementation: truth-preserving resume generation and contract-first cross-stack API changes.

## Roadmap

1. Candidate profile, experience, skill, and accomplishment CRUD.
2. Job-description ingestion and evidence-backed gap analysis.
3. Deterministic resume planning followed by constrained prose generation.
4. HTML/PDF rendering with snapshot versioning.
5. Application pipeline and interview-context retrieval.
6. Evaluation fixtures for factuality, relevance, and provider parity.

## Privacy and security

Candidate histories and job searches are sensitive. Resonance binds to loopback by default, stores data locally, keeps API keys server-side, and does not include telemetry by default. Before exposing the API to a network, add authentication, TLS, encrypted backups, and an explicit threat model.

## License

[MIT](LICENSE)

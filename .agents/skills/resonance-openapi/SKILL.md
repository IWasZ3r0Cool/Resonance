---
name: resonance-openapi
description: Add or change Resonance REST endpoints, OpenAPI schemas, or frontend API calls with synchronized contract, backend, generated TypeScript types, and verification.
---

# Resonance contract-first API changes

`api/openapi.json` is authoritative.

For each public API change:

1. Update the OpenAPI operation and reusable schemas, including stable `operationId` values and non-success responses.
2. Implement the Go handler and focused tests. Do not expose provider credentials, raw model internals, or candidate data that the operation does not require.
3. Run `pnpm --dir web generate:api` and use the generated types through `web/src/api`.
4. Update frontend loading, empty, success, and failure behavior where applicable.
5. Run `go test ./...`, `pnpm --dir web typecheck`, and `pnpm --dir web build`.

Treat generated `web/src/api/schema.d.ts` as an artifact: regenerate it rather than hand-editing it. Keep breaking changes explicit and prefer additive evolution until a versioned migration plan exists.

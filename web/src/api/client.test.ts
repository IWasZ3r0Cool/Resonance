import { afterEach, describe, expect, it, vi } from "vitest";

import { getRuntimeStatus } from "./client";

const health = { status: "ok", database: "ok", version: "test" } as const;
const capabilities = {
  agentSdk: "google-adk-go",
  database: "sqlite",
  llmConfigured: false,
  llmModel: "gpt-5-mini",
  llmProvider: "openai",
  localModelSupport: true,
} as const;

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("getRuntimeStatus", () => {
  it("combines the typed health and capabilities responses", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async (request: Request) => {
        const body = request.url.endsWith("/api/v1/health") ? health : capabilities;
        return Response.json(body);
      }),
    );

    await expect(getRuntimeStatus()).resolves.toEqual({ health, capabilities });
  });

  it("reports an unavailable health endpoint", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async (request: Request) => {
        if (request.url.endsWith("/api/v1/health")) {
          return Response.json({ error: "unavailable" }, { status: 503 });
        }
        return Response.json(capabilities);
      }),
    );

    await expect(getRuntimeStatus()).rejects.toThrow("local Resonance API is not available");
  });

  it("reports a missing capabilities response", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async (request: Request) => {
        if (request.url.endsWith("/api/v1/health")) {
          return Response.json(health);
        }
        return Response.json({ error: "unavailable" }, { status: 503 });
      }),
    );

    await expect(getRuntimeStatus()).rejects.toThrow("runtime capabilities");
  });
});

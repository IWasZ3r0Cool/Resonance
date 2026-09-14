import createClient from "openapi-fetch";

import type { components, paths } from "./schema";

export type Health = components["schemas"]["HealthResponse"];
export type Capabilities = components["schemas"]["CapabilitiesResponse"];

export async function getRuntimeStatus(signal?: AbortSignal): Promise<{
  health: Health;
  capabilities: Capabilities;
}> {
  const client = createClient<paths>({
    baseUrl: import.meta.env.VITE_API_BASE_URL || window.location.origin,
    fetch: globalThis.fetch,
  });
  const [healthResult, capabilitiesResult] = await Promise.all([
    client.GET("/api/v1/health", { signal }),
    client.GET("/api/v1/meta/capabilities", { signal }),
  ]);

  if (healthResult.error || !healthResult.data) {
    throw new Error("The local Resonance API is not available.");
  }
  if (capabilitiesResult.error || !capabilitiesResult.data) {
    throw new Error("The API did not return its runtime capabilities.");
  }

  return {
    health: healthResult.data,
    capabilities: capabilitiesResult.data,
  };
}

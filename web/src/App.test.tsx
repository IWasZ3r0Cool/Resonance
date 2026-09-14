import { render, screen } from "@testing-library/react";
import { ThemeProvider } from "@mui/material";
import { beforeEach, describe, expect, it, vi } from "vitest";

import App from "./App";
import { getRuntimeStatus } from "./api/client";
import { theme } from "./theme";

vi.mock("./api/client", () => ({
  getRuntimeStatus: vi.fn(),
}));

const mockedGetRuntimeStatus = vi.mocked(getRuntimeStatus);

function renderApp() {
  return render(
    <ThemeProvider theme={theme}>
      <App />
    </ThemeProvider>,
  );
}

describe("App", () => {
  beforeEach(() => {
    mockedGetRuntimeStatus.mockReset();
  });

  it("shows the configured local runtime", async () => {
    mockedGetRuntimeStatus.mockResolvedValue({
      health: { status: "ok", database: "ok", version: "test" },
      capabilities: {
        agentSdk: "google-adk-go",
        database: "sqlite",
        llmConfigured: true,
        llmModel: "qwen3:8b",
        llmProvider: "ollama",
        localModelSupport: true,
      },
    });

    renderApp();

    expect(await screen.findByText("Local API ready")).toBeInTheDocument();
    expect(screen.getByText("ollama · qwen3:8b")).toBeInTheDocument();
    expect(screen.getByText("configured")).toBeInTheDocument();
  });

  it("explains how to start an unavailable API", async () => {
    mockedGetRuntimeStatus.mockRejectedValue(new Error("The local Resonance API is not available."));

    renderApp();

    expect(await screen.findByText("API offline")).toBeInTheDocument();
    expect(screen.getByText(/go run \.\/cmd\/api/)).toBeInTheDocument();
  });
});

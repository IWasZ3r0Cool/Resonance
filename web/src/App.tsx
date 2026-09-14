import AutoAwesomeRounded from "@mui/icons-material/AutoAwesomeRounded";
import CheckCircleRounded from "@mui/icons-material/CheckCircleRounded";
import DataObjectRounded from "@mui/icons-material/DataObjectRounded";
import DescriptionRounded from "@mui/icons-material/DescriptionRounded";
import HistoryRounded from "@mui/icons-material/HistoryRounded";
import LockRounded from "@mui/icons-material/LockRounded";
import WorkOutlineRounded from "@mui/icons-material/WorkOutlineRounded";
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Container,
  Stack,
  Typography,
} from "@mui/material";
import type { SvgIconComponent } from "@mui/icons-material";
import { useEffect, useState } from "react";

import { getRuntimeStatus } from "./api/client";
import type { Capabilities, Health } from "./api/client";

type RuntimeState =
  | { kind: "loading" }
  | { kind: "error"; message: string }
  | { kind: "ready"; health: Health; capabilities: Capabilities };

interface Feature {
  icon: SvgIconComponent;
  eyebrow: string;
  title: string;
  description: string;
}

const features: Feature[] = [
  {
    icon: DataObjectRounded,
    eyebrow: "Source of truth",
    title: "Capture the whole story",
    description:
      "Keep detailed roles, reusable achievement bullets, skills, hats, and evidence in one private workspace.",
  },
  {
    icon: AutoAwesomeRounded,
    eyebrow: "Evidence-backed AI",
    title: "Tailor without fiction",
    description:
      "Match a job to what you can prove, surface real gaps, and draft from supported facts only.",
  },
  {
    icon: HistoryRounded,
    eyebrow: "Application memory",
    title: "Know exactly what you sent",
    description:
      "Retain every job description, resume version, and application decision for interview day.",
  },
];

function App() {
  const [runtime, setRuntime] = useState<RuntimeState>({ kind: "loading" });

  useEffect(() => {
    const controller = new AbortController();
    getRuntimeStatus(controller.signal)
      .then(({ health, capabilities }) => {
        setRuntime({ kind: "ready", health, capabilities });
      })
      .catch((error: unknown) => {
        if (!controller.signal.aborted) {
          setRuntime({
            kind: "error",
            message: error instanceof Error ? error.message : "The local API is unavailable.",
          });
        }
      });
    return () => controller.abort();
  }, []);

  return (
    <Box component="main" className="min-h-screen overflow-hidden">
      <Container maxWidth="lg" className="relative py-6 sm:py-10">
        <Stack
          component="header"
          direction="row"
          sx={{ alignItems: "center", justifyContent: "space-between" }}
          className="mb-16 sm:mb-24"
        >
          <Stack direction="row" spacing={1.5} sx={{ alignItems: "center" }}>
            <Box
              aria-hidden="true"
              className="grid h-10 w-10 place-items-center rounded-2xl bg-[#1f1d2b] text-white"
            >
              <DescriptionRounded fontSize="small" />
            </Box>
            <Box>
              <Typography sx={{ fontWeight: 800, lineHeight: 1.05 }}>
                Resonance
              </Typography>
              <Typography variant="caption" color="text.secondary">
                Resume Engine
              </Typography>
            </Box>
          </Stack>
          <RuntimeBadge runtime={runtime} />
        </Stack>

        <section className="grid items-end gap-10 lg:grid-cols-[1.35fr_0.65fr]">
          <Box>
            <Chip
              icon={<LockRounded />}
              label="Local-first · Your evidence stays yours"
              variant="outlined"
              className="mb-7 bg-white/50 backdrop-blur"
            />
            <Typography variant="h1" component="h1" sx={{ maxWidth: 820 }}>
              Every version of you,
              <Box component="span" sx={{ color: "primary.main" }}>
                {" "}in tune with the role.
              </Box>
            </Typography>
            <Typography
              variant="h6"
              color="text.secondary"
              sx={{ fontWeight: 400, lineHeight: 1.65, maxWidth: 720 }}
              className="mt-7"
            >
              Build one complete career record. Resonance finds the strongest truthful signal for
              each opportunity, explains what is missing, and remembers what you applied with.
            </Typography>
            <Stack direction={{ xs: "column", sm: "row" }} spacing={1.5} className="mt-9">
              <Button variant="contained" size="large" startIcon={<WorkOutlineRounded />} disabled>
                Add your first experience
              </Button>
              <Button variant="text" size="large" href="/openapi.json">
                Explore the API contract
              </Button>
            </Stack>
          </Box>

          <StatusCard runtime={runtime} />
        </section>

        <section aria-labelledby="foundation-heading" className="mt-24 sm:mt-32">
          <Typography id="foundation-heading" variant="overline" color="primary" sx={{ fontWeight: 800 }}>
            The foundation
          </Typography>
          <Typography variant="h2" component="h2" className="mt-2 max-w-2xl">
            Designed around evidence, not keyword theater.
          </Typography>
          <div className="mt-9 grid gap-4 md:grid-cols-3">
            {features.map((feature, index) => (
              <FeatureCard key={feature.title} feature={feature} index={index} />
            ))}
          </div>
        </section>

        <Box component="footer" className="mt-20 border-t border-black/10 py-8">
          <Typography variant="body2" color="text.secondary">
            Open source under the MIT License · Foundation release
          </Typography>
        </Box>
      </Container>
    </Box>
  );
}

function RuntimeBadge({ runtime }: { runtime: RuntimeState }) {
  if (runtime.kind === "loading") {
    return <Chip size="small" icon={<CircularProgress size={14} />} label="Checking local API" />;
  }
  if (runtime.kind === "error") {
    return <Chip size="small" color="warning" variant="outlined" label="API offline" />;
  }
  return <Chip size="small" color="success" icon={<CheckCircleRounded />} label="Local API ready" />;
}

function StatusCard({ runtime }: { runtime: RuntimeState }) {
  return (
    <Card className="bg-[#1f1d2b] text-white">
      <CardContent className="p-6 sm:p-7">
        <Typography variant="overline" sx={{ color: "rgba(255,255,255,.56)" }}>
          Runtime
        </Typography>
        {runtime.kind === "loading" && (
          <Stack direction="row" spacing={2} sx={{ alignItems: "center" }} className="mt-5">
            <CircularProgress size={22} color="inherit" />
            <Typography>Connecting to the local API…</Typography>
          </Stack>
        )}
        {runtime.kind === "error" && (
          <Alert severity="warning" className="mt-4">
            {runtime.message} Start it with <code>go run ./cmd/api</code>.
          </Alert>
        )}
        {runtime.kind === "ready" && (
          <Stack spacing={2.25} className="mt-5">
            <RuntimeRow label="API" value={`v${runtime.health.version}`} />
            <RuntimeRow label="Database" value="SQLite · ready" />
            <RuntimeRow
              label="Model"
              value={`${runtime.capabilities.llmProvider} · ${runtime.capabilities.llmModel}`}
            />
            <RuntimeRow
              label="Generation"
              value={runtime.capabilities.llmConfigured ? "configured" : "add a provider key"}
            />
          </Stack>
        )}
      </CardContent>
    </Card>
  );
}

function RuntimeRow({ label, value }: { label: string; value: string }) {
  return (
    <Stack direction="row" spacing={2} sx={{ justifyContent: "space-between" }}>
      <Typography sx={{ color: "rgba(255,255,255,.56)" }}>{label}</Typography>
      <Typography sx={{ textAlign: "right", fontWeight: 700 }}>
        {value}
      </Typography>
    </Stack>
  );
}

function FeatureCard({ feature, index }: { feature: Feature; index: number }) {
  const Icon = feature.icon;
  return (
    <Card className="h-full bg-white/70 backdrop-blur">
      <CardContent className="p-6 sm:p-7">
        <Stack direction="row" sx={{ justifyContent: "space-between", alignItems: "flex-start" }}>
          <Box className="grid h-11 w-11 place-items-center rounded-2xl bg-[#eee9ff] text-[#5b3fd1]">
            <Icon />
          </Box>
          <Typography variant="caption" color="text.secondary" sx={{ fontWeight: 800 }}>
            0{index + 1}
          </Typography>
        </Stack>
        <Typography variant="overline" color="secondary" sx={{ fontWeight: 800 }} className="mt-6 block">
          {feature.eyebrow}
        </Typography>
        <Typography variant="h5" sx={{ fontWeight: 750 }} className="mt-1">
          {feature.title}
        </Typography>
        <Typography color="text.secondary" sx={{ lineHeight: 1.7 }} className="mt-3">
          {feature.description}
        </Typography>
      </CardContent>
    </Card>
  );
}

export default App;

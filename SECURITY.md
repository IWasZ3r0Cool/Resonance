# Security policy

## Supported versions

Resonance is pre-release software. Security fixes are applied to the latest commit on `main`; no older release line is supported yet.

## Reporting a vulnerability

Please do not open a public issue for a suspected vulnerability or include real candidate data, resumes, job-search history, credentials, or exploit details in public discussions.

Use GitHub's private vulnerability reporting feature for this repository. Include:

- the affected commit or version;
- a concise reproduction using synthetic data;
- the impact and required attacker access;
- any suggested mitigation, if known.

Maintainers should acknowledge a report within seven days, keep the reporter informed while it is assessed, and coordinate disclosure after a fix is available.

## Security model

The development server is intended for loopback access only. It does not yet provide authentication or TLS and must not be exposed directly to a public or untrusted network.

Provider credentials belong only in server environment variables. Candidate data and generated resume artifacts are sensitive local data. Logs, tests, bug reports, and telemetry must use synthetic content and must not capture either category.

## Automated checks

Pull requests run tests, race detection, static analysis, contract-drift detection, dependency review, CodeQL, Go vulnerability scanning, and npm's high-severity audit. These checks reduce risk but do not replace threat modeling and human review.

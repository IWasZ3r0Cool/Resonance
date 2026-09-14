---
name: resonance-domain
description: Implement or review Resonance candidate records, job-gap analysis, resume generation, application tracking, or related persistence while preserving factual provenance and reproducibility.
---

# Resonance domain work

Treat the candidate record as evidence, not creative source material.

## Required outcomes

- Trace every generated resume claim to one or more stored candidate facts.
- Distinguish stated facts from model inference. Never promote an inference into the source of truth without candidate confirmation.
- Store reusable accomplishment text as `resume_bullets`. Use `bullet_skills` for demonstrated technologies/capabilities and typed `bullet_tags` for hats or differentiators; retain association provenance and rationale.
- Keep every experience, bullet, skill, tag, and association within one candidate profile. Reject cross-candidate evidence links at the persistence boundary.
- When a job requirement lacks evidence, return a visible gap or clarification prompt. Reordering, omitting, or carefully reframing true facts is allowed; fabrication and unjustified strengthening are not.
- Preserve the job-description snapshot, selected evidence IDs, prompt/workflow version, provider/model identifier, and output for any resume version linked to an application.
- Allow resume drafts to change until finalization. Final versions are immutable, applications may reference only final versions, and a revision creates a new draft linked to the finalized same-candidate parent.
- Avoid logging resume content, candidate PII, job-search activity, or credentials by default.

## Implementation guidance

Keep deterministic selection and validation separate from prose generation. Prefer structured model output with evidence identifiers, then reject claims whose identifiers do not resolve to eligible source records. Test the validator without a live model.

When changing storage, add a forward-only migration. When changing HTTP behavior, also invoke the `resonance-openapi` skill.

Before finishing, test at least one complete-evidence case, one genuine-gap case, and one attempted unsupported-claim case for the behavior being changed.

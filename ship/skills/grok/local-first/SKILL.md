---
name: local-first
description: >-
  Route work through the cheapest capable tier: scripts and CLIs first, then a
  local model (Ollama), then an agreed cloud fallback (DeepSeek by default),
  and report the tokens each tier spent. Use when the user says local models,
  Ollama, GPU, DeepSeek, minimize tokens, or "don't burn tokens", or when a
  loop such as scanner autofix needs an LLM.
user-invocable: true
---

# local-first

## Pstack root

- `principle-laziness-protocol`: the least machinery that does the job.
  Mechanical work gets a script, not a model.
- `principle-build-the-lever`: when a mechanical job repeats, write the tool
  that does it.
- `principle-guard-the-context-window`: structured tool output over model
  prose.

Wade-specific parameters (treat as configuration, not doctrine): laptop
Ollama, the "habitat" Mac mini Qwen-class host (not live; mocked), and
DeepSeek as the cloud fallback. In this repo the cascade is implemented by
`frontier ai route|complete|tokens`.

## Cascade (least tokens first)

| Tier | What | Use when | Fall through when |
|---|---|---|---|
| 0 | Process: Frontier commands, scripts, scanners, `gh`, `gcloud`, tests | The job is mechanical: version bumps, Dockerfile `USER`/`HEALTHCHECK`, pinning Actions to SHAs, redacting a known secret | No deterministic tool can finish it |
| 1 | Laptop Ollama, already-pulled coding model | An LLM is needed and the probe answers | One probe fails or the model is too weak for the task |
| 2 | Habitat / Mac mini Qwen-class (~27B) | Host is live | Always, until the host exists: mock as unavailable and record `0 (mocked fallthrough)` |
| 3 | DeepSeek API (`DEEPSEEK_API_KEY` or `frontier ai secrets`) | Every local tier is down or insufficient | Not configured: stop and report |

Do not use any other cloud coding model for autofix unless the user names one.
Move back up the cascade as soon as a higher tier is healthy.

## Steps

1. **Classify the job.** If tier 0 can do it, do it and skip the rest.
2. **Probe once.** `frontier ai route` (or one request to the local chat
   endpoint, e.g. `curl -fsS --max-time 5 $OLLAMA_HOST/api/tags`). On failure,
   stop using local subagents for this session. Do not retry in a loop.
3. **Run** on the first healthy tier. Keep prompts short; ask for patches or
   structured output, not prose.
4. **Do not poll** for GPU or job status. Report the current state once when
   asked.
5. **Report spend.** For PRs, paste the token report from
   `frontier ai tokens` (reads `.frontier/tokens.jsonl`). Format:
   [token report](references/token-report.md).

## Done when

- [ ] The report names the tier that did the work and why higher-preference
      tiers were skipped (probe result or "mechanical, tier 0").
- [ ] Local availability was probed at most once per session (plus once after
      a known recovery).
- [ ] PRs produced through the cascade carry the token report table, with
      tier 2 recorded as `0 (mocked fallthrough)` while the host is not live.

## Conflicts

- Grok Bot or the host agent itself is the orchestrator; this skill governs
  *which model does delegated LLM work*, not whether the orchestrator thinks.
- Scanner autofix uses this cascade; see `scanner-hotfix` for the loop.
- Secrets: never print or commit the DeepSeek key.
  `frontier ai secrets status` shows where it would load from without the
  value.

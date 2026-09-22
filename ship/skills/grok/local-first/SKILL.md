---
name: local-first
description: >-
  Prefer local Ollama when it is available; otherwise DeepSeek API is an allowed cloud fallback.
  Prefer tools/skills/process over any LLM when the job is mechanical. Use when the user says
  local models, Ollama, GPU, minimize tokens, DeepSeek, or "don't burn tokens."
user-invocable: true
---

# local-first

## Priority order (least tokens first)

1. **Process first.** If Frontier skills, scripts, scanners, `gh`, `gcloud`, or tests can finish the job, do that. Do not spawn an LLM subagent for mechanical work (version bumps, Dockerfile USER/HEALTHCHECK, pinning Actions to SHAs, redacting known secrets).
2. **Laptop local (Ollama).** Prefer already-pulled models on the developer GPU. Treat a weak laptop GPU as an **auto-fallthrough** — one light probe, then continue down the cascade; do not thrash.
3. **Habitat / Mac mini Qwen-class.** Mid-tier on-prem coding model (e.g. community-tuned ~27B). Until that host is live, **mock** the probe as unavailable and fall through; still record the route as `0 (mocked fallthrough)` in PR token reports so the workflow maps future local capacity.
4. **Frontier AI = DeepSeek API.** When every local tier is down/insufficient and an LLM is still required, use DeepSeek (respect existing env keys / project config). Do **not** use other cloud coding models for autofix unless the user names one.
5. Stay ready to switch back up the cascade the moment a higher tier is healthy.

## Steps

1. Probe local availability lightly (e.g. one request to the local chat endpoint). On connection failure, pause local-subagent use for this session and fall back to DeepSeek or pure tools.
2. Prefer already-pulled local coding models and the local GPU when present.
3. Minimize tokens everywhere. Prefer structured tool output over long model prose.
4. Do not poll for GPU or job status. Report current artifacts once when asked.
5. Scanner autofix and similar loops: runner opens issues; a **local** (or DeepSeek) agent may open `fix/vuln-*` PRs into `dev` under human control — see skill `scanner-hotfix`.

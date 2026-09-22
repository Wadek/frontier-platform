---
name: local-first
description: >-
  Prefer local Ollama when it is available; otherwise DeepSeek API is an allowed cloud fallback.
  Prefer tools/skills/process over any LLM when the job is mechanical. Use when the user says
  local models, Ollama, GPU, minimize tokens, DeepSeek, or "don't burn tokens."
user-invocable: true
---

# local-first

## Priority order

1. **Process first.** If Frontier skills, scripts, scanners, `gh`, `gcloud`, or tests can finish the job, do that. Do not spawn an LLM subagent for mechanical work.
2. **Local agents when available.** If Ollama (or another local OpenAI-compatible endpoint) is up, use it for drafting/planning. Do not assume Ollama is always running.
3. **DeepSeek API fallback.** When local agents are unavailable and an LLM is still required, use DeepSeek via its API (respect existing env keys / project config). Stay ready to switch back to local the moment it is healthy.
4. Avoid other cloud LLM spend unless the user names a provider.

## Steps

1. Probe local availability lightly (e.g. one request to the local chat endpoint). On connection failure, pause local-subagent use for this session and fall back to DeepSeek or pure tools.
2. Prefer already-pulled local coding models and the local GPU when present.
3. Minimize tokens everywhere. Prefer structured tool output over long model prose.
4. Do not poll for GPU or job status. Report current artifacts once when asked.
5. Scanner autofix and similar loops: runner opens issues; a **local** (or DeepSeek) agent may open `fix/vuln-*` PRs into `dev` under human control — see skill `scanner-hotfix`.

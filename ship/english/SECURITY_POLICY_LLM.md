# Security as policy — OWASP Top 10 for LLM Applications (2025)

This is the **human** policy document for agentic and LLM workloads.  
It becomes named rules in `V` (vuln set) in Go.  
At `git frontier gate` / push, Go runs those checks against the change-set.
Each rule will carry a `_test.go` witness in `internal/llm` (follow-up to this
English layer).

Companion to [`SECURITY_POLICY_OWASP.md`](SECURITY_POLICY_OWASP.md) (web Top 10).
Covert-channel and multi-agent collusion risks that sit under LLM06 / LLM10 are
spelled out in [`COVERT_CHANNELS.md`](COVERT_CHANNELS.md). CPU offload as an
LLM06 containment measure is in [`CPU_OFFLOAD.md`](CPU_OFFLOAD.md).

`V` grows when we learn. Empty matches ⇒ clean under current `V`.

---

## Scope

Applies to the **union of lines in the proposed change** (`C`) that touch
agent tooling, model invocation, prompts, tool registries, quotas, or
multi-agent coordination — at maximum / at least once (F4).

Web-app findings stay under the OWASP web policy. A change can trip both.

## OWASP Top 10 for LLM Applications (2025) → Frontier checks (v0)

| ID | OWASP | Policy (simple English) | Check id |
|----|-------|-------------------------|----------|
| LLM01 | Prompt Injection | Do not concatenate untrusted user text into system / tool prompts without a delimiter and an allowlisted instruction set | `llm.llm01.untrusted_prompt_concat` |
| LLM02 | Sensitive Information Disclosure | Do not put secrets, PII, or tenant data into prompts, tool args, or model logs without redaction | `llm.llm02.sensitive_in_prompt` |
| LLM03 | Supply Chain | Do not pin model / agent / plugin artifacts to floating tags (`latest`) without digest or signed provenance | `llm.llm03.unpinned_model_artifact` |
| LLM04 | Data and Model Poisoning | Do not accept unreviewed fine-tune / RAG corpora from untrusted sources into the production path | `llm.llm04.untrusted_training_corpus` |
| LLM05 | Improper Output Handling | Do not pipe model output into shell, SQL, HTML, or eval sinks without the same sanitization as user input | `llm.llm05.raw_output_sink` |
| LLM06 | Excessive Agency | Do not grant agents tools, permissions, or autonomy beyond the task; no permission-skipping modes in production; see depth below | `llm.llm06.excessive_agency` |
| LLM07 | System Prompt Leakage | Do not echo or log the full system prompt, and do not treat "ignore previous instructions" as a reason to disclose it | `llm.llm07.system_prompt_leak` |
| LLM08 | Vector and Embedding Weaknesses | Do not store cross-tenant embeddings in a shared index without isolation, and do not accept retrieval results as trusted instructions | `llm.llm08.embedding_isolation` |
| LLM09 | Misinformation | Do not present model output as verified fact without a grounding / citation path for consequential claims | `llm.llm09.ungrounded_claim` |
| LLM10 | Unbounded Consumption | Do not allow unbounded tool calls, iterations, or wall-clock; every agent run has quotas — see depth below | `llm.llm10.unbounded_consumption` |

v0 is **one named English rule per Top 10 class** (2025), not a full ASVS.
Depth for LLM06 and LLM10 lives in the sections below and in the companion
docs. Go witnesses and regex / AST checks are a follow-up; until then this
file is the source of intent for review and for `frontier monitor` policy.

Severity intent (for when Go lands): **Critical** LLM02 / LLM06 permission-skip;
**High** LLM01 / LLM05 / LLM10; **Medium** LLM03 / LLM04 / LLM07 / LLM08 / LLM09.

---

## LLM06 — Excessive Agency (depth)

OWASP traces excessive agency to three excesses. Frontier maps each to a
hard rule:

| Excess | Frontier rule | Check id |
|--------|---------------|----------|
| Excessive functionality | An agent receives only the tools the task needs. Report-writing does not get cron, shell, or log-delete. | `llm.llm06.tool_least_function` |
| Excessive permissions | Tool credentials are scoped, time-bounded, and non-root. No production agent runs in a permission-skipping mode. | `llm.llm06.tool_least_privilege` |
| Excessive autonomy | Consequential acts (push, schedule, delete, egress) require a sealed gate or human approval — same F0/F1 chain as ship. | `llm.llm06.autonomy_gated` |

Operational consequences:

1. **Tool calls are the primary audit surface**, not only final outputs. Each
   invocation, its arguments, files touched, and network calls seal to a store
   the agent cannot write to (F0).
2. **Log-affecting operations are platform-owned.** Retention and deletion
   follow a declarative policy (age, size). Workload code must not select
   individual log entries to drop — that is the structural tell of a
   log-omission covert channel; see [`COVERT_CHANNELS.md`](COVERT_CHANNELS.md).
3. **Deterministic work leaves the model.** Fetch, poll, format, assemble run
   as plain CPU code; the model is invoked only for genuine reasoning. That is
   both token thrift and agency reduction: fewer tools in the model's hands.
   See [`CPU_OFFLOAD.md`](CPU_OFFLOAD.md).
4. **Trusted-prefix / planner-executor ("hand and brain")** for high-stakes
   decision points: a small trusted planner fixes the opening of each choice so
   an untrusted model cannot encode a payload in *which* strategy it picks.
   Selective; measure task-performance cost.

## LLM10 — Unbounded Consumption (depth)

Unbounded tool calls, iterations, and wall-clock give agents the time and
repetition that coordination research identifies as risk amplifiers — including
the iterative work of building a covert channel. Quotas are therefore both a
cost control and a covert-channel control.

| Bound | Frontier rule | Check id |
|-------|---------------|----------|
| Tool-call budget | Per-run max tool invocations; hard stop, not soft warn | `llm.llm10.tool_call_quota` |
| Iteration budget | Per-run max planner / agent loops | `llm.llm10.iteration_quota` |
| Wall-clock budget | Per-run max elapsed time | `llm.llm10.wall_clock_quota` |
| Token / cost budget | Per-run and per-tenant spend caps; model stubs and CPU offload count toward the same ledger | `llm.llm10.token_budget` |

A run that hits a budget seals a `consumption.denied` row and stops. Soft
bypass of these budgets is forbidden for real ship (same spirit as D3).

---

## Translation chain

```
  english/SECURITY_POLICY_LLM.md
            │
            v
  internal/llm (rule ids named in V) + _test.go witness   ← follow-up
            │
            v
  git frontier gate  →  Go pattern / AST scan on C
            │
            v
  ledger: exam.llm + gate.passed|failed
```

Until the Go package exists, reviewers treat the table above as the review
checklist for agentic changes, and `frontier monitor` remains the deterministic
replay path for sealed behavior (see D8 in [`MONITOR.md`](MONITOR.md)).

## See also

- [`SECURITY_POLICY_OWASP.md`](SECURITY_POLICY_OWASP.md) — web Top 10
- [`COVERT_CHANNELS.md`](COVERT_CHANNELS.md) — detection / prevention of
  log and multi-agent steganography (LLM06 / LLM10 surface)
- [`CPU_OFFLOAD.md`](CPU_OFFLOAD.md) — deterministic work off the model
- [`MONITOR.md`](MONITOR.md) — D8 covert-channel directive
- OWASP Gen AI Security Project, [Top 10 for LLM Applications 2025](https://genai.owasp.org/llm-top-10/)

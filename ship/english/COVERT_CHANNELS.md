# Covert channels — detection and prevention (policy)

This is the **human** policy for detecting and preventing covert /
steganographic channels in logs and multi-agent coordination. It is
**detection- and prevention-only**. It does not describe how to build a
channel, does not map a message onto gaps or timestamps, and does not ship
an encoder or decoder.

English states the policy; Go witnesses under `internal/monitor` (and a
future `internal/covert`) are the follow-up machine check. Until then,
reviewers and `frontier monitor` treat this file as source of intent.

Maps onto [`SECURITY_POLICY_LLM.md`](SECURITY_POLICY_LLM.md) **LLM06**
(Excessive Agency) and **LLM10** (Unbounded Consumption), and onto monitor
directive **D8** in [`MONITOR.md`](MONITOR.md).

---

## What a covert channel is here

A **covert channel** is a communication path the system's designers never
intended and its monitors were not built to see (Lampson, 1973). Three
families matter for Frontier:

| Family | Signal sits in | Typical cover |
|--------|----------------|---------------|
| Storage | Values the system stores or transmits | Unused fields, log-message content |
| Timing | When events occur | Inter-arrival times, delays, cadence |
| Event-sequence | Order / selection / **absence** of events | Which heartbeats, audits, or job outputs are missing |

**Steganography** aims to hide that a message exists, not only its content.
A successful scheme produces output a content monitor cannot reliably tell
apart from ordinary output.

Logs are an attractive carrier: they are often plaintext, event-triggered
(an outsider can cause entries without host access), asynchronously
retrievable, and persisted in several places (Lamshöft et al., 2022).

---

## The structural tell (what to look for)

Content-based monitors miss channels whose signal is in **structure**, not
text. The tell Frontier cares about:

> A "maintenance", "retention", or "filter" job whose **excluded /
> omitted / deleted set is an arbitrary information-carrying list**, not a
> rule (age, size, severity, allowlist).

Concrete review questions — all answerable without building a channel:

1. **Is the omission rule declarative and platform-owned?** Age-based
   rotation, size caps, severity floors are rules. A hand-picked list of
   timestamps, sequence numbers, or entry ids is not.
2. **Could a different operator reconstruct the same exclusions from
   policy alone?** If no, the set is carrying information beyond policy.
3. **Do independent log sources for the same activity agree?** A gap in
   one source while others keep recording is a tamper / channel indicator
   (gap-and-burst analysis; IRFlow Timeline).
4. **Does a regular stream's cadence deviate?** Heartbeats, health checks,
   and scheduled jobs have predictable rhythms. Irregular gap patterns and
   unusual distributions of numeric fields (for example reported delays)
   are `watch` signals even when every surviving line is authentic.

The same structural lesson applies to agents: the payload may sit in
*what the agent chooses to do* (strategy, algorithm, narrative), not in
*how it writes*. Paraphrasing leaves those payloads in place
(Krzyzanowski et al., 2026).

Illustrative only — no meaning is encoded:

```
10:00:00  INFO  heartbeat ok
10:00:10  INFO  heartbeat ok
          [entry removed by maintenance job]
10:00:30  INFO  heartbeat ok
          [entry removed by maintenance job]
10:01:00  INFO  heartbeat ok
```

Every surviving line is valid. Anomaly detection requires looking at the
**timeline structure**, not the text.

---

## How `frontier monitor` flags it (deterministic, no model)

`frontier monitor` is a deterministic replay of sealed ledger rows — not a
model, no tokens ([`MONITOR.md`](MONITOR.md)). Directive **D8** adds the
covert-channel check to that replay:

| Signal | Verdict | Notes |
|--------|---------|-------|
| Omitted / excluded set on a monitored job is not derivable from a declared rule | `violation` | Arbitrary list = information-carrying |
| Anomalous gap-and-burst timing vs the job's own baseline | `watch` (promote to `violation` when correlated across sources) | Advise-only until Go witness lands |
| Retention / logging-config / maintenance-job change without an independent audit seal | `watch` | "Audit the auditors" |
| Ledger hash chain broken around a gap | `tampered` (D0) | Deletion that also breaks F0 |

v0 scope for D8 is **sealed job metadata and ledger timing**, not full
syslog forensics. Cross-source gap analysis against node / audit / flow
logs is an operator procedure that feeds findings back into the ledger;
the monitor then replays those seals.

---

## Defenses (capability, bandwidth, coordination)

Perfect detection of information-theoretically secure schemes is not
achievable. Frontier's stance: **limit capability, bandwidth, and
coordination opportunity, and make tampering evident.**

### 1. Least privilege and tool-use monitoring (LLM06)

- Tool calls are the primary audit surface. Seal each invocation off-host
  where the agent cannot write.
- No production agent in a permission-skipping mode.
- Capabilities per task and time-bounded. A report writer does not get
  cron. A maintenance job does not choose individual entries to delete.
- Log retention / deletion is a **platform-owned declarative policy**,
  never workload-selectable per entry.

### 2. Trusted-prefix / planner-executor ("hand and brain")

Where an untrusted, highly capable model makes consequential choices
(strategy, plan, code-path), a small **trusted** planner writes the first
*k* tokens (or otherwise fixes the opening) at each decision point so the
choice cannot depend on a hidden payload (Krzyzanowski et al., 2026).
Apply selectively; measure task-performance cost. See also
[`SECURITY_POLICY_LLM.md`](SECURITY_POLICY_LLM.md) LLM06 depth.

### 3. Bound consumption (LLM10)

Per-agent quotas on tool calls, iterations, and wall-clock limit the
iterative work needed to build and tune a channel. Same budgets as
[`SECURITY_POLICY_LLM.md`](SECURITY_POLICY_LLM.md) LLM10.

### 4. Egress, east-west, and shared-artefact controls

- Default-deny egress; allowlist trusted endpoints only (Motwani et al.).
- Restrict east-west traffic so pods talk only to required services.
- Separate tenants and agents at the storage layer; shared writable
  volumes between trust zones are covert storage channels.
- Block cloud instance-metadata endpoints from workload pods.

### 5. Log integrity

- Forward logs off-host promptly to append-only or write-once storage.
- Chain or sign entries; keep sequence numbers so a missing record is
  **provably** missing.
- Baseline regular streams; alert on cadence and gap-pattern deviation.
- Changes to retention, logging config, or maintenance jobs seal to an
  independent store and are themselves reviewed.

### 6. Reduce coordination opportunity

Tool use has largely removed the *capability* barrier to undetectable
steganography; **coordination** (shared scheme, key, parameters) remains
the bottleneck (Rippin et al., 2026). Reduce shared focal points:

- Limit repeated unmonitored interaction between the same agent pairs.
- Reset shared state between tasks where possible.
- Insert monitors between agents.
- Prefer per-instance diversity of prompts and environments over a single
  shared configuration.

### 7. CPU offload as containment

Deterministic agent steps (fetch, poll, format, assemble) run as plain
CPU code; the model is invoked only for genuine reasoning. Fewer tools in
the model's hands shrinks the covert-channel surface and the token bill.
See [`CPU_OFFLOAD.md`](CPU_OFFLOAD.md).

---

## What this policy does not do

- Does not provide an encoder, decoder, codebook, or gap↔symbol mapping.
- Does not claim complete detection of secure steganographic schemes.
- Does not replace content monitoring — content checks remain necessary;
  they are not sufficient.
- Does not authorize permission-skipping modes "to research the threat."

---

## References

1. Lampson, B. W. (1973). A note on the confinement problem. *Communications of the ACM*, 16(10), 613–615. https://doi.org/10.1145/362375.362389
2. Lamshöft, K., Neubert, T., Hielscher, J., Vielhauer, C., & Dittmann, J. (2022). *Knock, Knock, Log: Threat Analysis, Detection & Mitigation of Covert Channels in Syslog using Port Scans as Cover.* DFRWS EU 2022; Forensic Science International: Digital Investigation, 40, 301335. https://doi.org/10.1016/j.fsidi.2022.301335
3. Motwani, S. R., et al. (2024). *Secret Collusion among AI Agents: Multi-Agent Deception via Steganography.* NeurIPS 2024. https://arxiv.org/html/2402.07510v4
4. Rippin, J. L., Marshall, S. C., Africa, D. D., & Schroeder de Witt, C. (2026). *Tool Use Enables Undetectable Steganography in Multi-Agent LLM Systems.* arXiv:2606.28425. https://arxiv.org/html/2606.28425v1
5. Krzyzanowski, R., Arcuschin, I., Lee, M., Meyer, B., & Lange, G. (2026). *Hand and Brain: Defenses against Agentic Steganography in Language Models.* ICML 2026. https://icml.cc/virtual/2026/76954
6. OWASP Gen AI Security Project. *LLM06:2025 Excessive Agency.* https://genai.owasp.org/llmrisk/llm062025-excessive-agency/
7. OWASP Gen AI Security Project. *LLM10:2025 Unbounded Consumption.* https://genai.owasp.org/llmrisk/llm102025-unbounded-consumption/
8. IRFlow Timeline. *Gap & Burst Analysis.* https://r3nzsec.github.io/irflow-timeline/features/gap-burst-analysis

## See also

- [`SECURITY_POLICY_LLM.md`](SECURITY_POLICY_LLM.md)
- [`SECURITY_POLICY_OWASP.md`](SECURITY_POLICY_OWASP.md) — A09 Logging/Monitoring Failures
- [`MONITOR.md`](MONITOR.md) — directive D8
- [`CPU_OFFLOAD.md`](CPU_OFFLOAD.md)
- [`H_HYGIENE.md`](H_HYGIENE.md) — provenance marks (related audit surface)

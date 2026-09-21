# Control — doctrine (frontier-control)

frontier-control coordinates agents working on projects. It routes, clears, transfers, completes, and debriefs. It never does the work itself, and a model never decides which model to use.

## Components

| Component | Job |
|---|---|
| taxonomist | manage per-project taxonomy branches (`control taxonomy …`) |
| data index | registry of projects + agents (`control roster`), derived from the file tree |
| pricing | peak/off-peak gating + daily pricing-page check (`control pricing`) |
| job service | scheduled one-shots via the OS scheduler — never polling (`control ops jobs`) |
| workflow engine | FSM over an append-only journal; every state rehydrates fresh context |
| review | review policy: scope, pricing, egress, deny > ask > allow (`control workflow review`) |
| providers | thin harness adapters: same workflow, `--provider local\|deepseek-flash\|deepseek-v4-pro\|dsh` |
| monitor | replay the ledger against directives (`control monitor`) |
| prove | Proven/Theory harness — no claim without a passing test (`control prove`) |

## Workflow FSM

```
queued → review → preparing (hydrate own taxonomy) → running → completed → debriefed → closed
  review → deferred   when an expensive provider is requested at peak (job service reschedules off-peak)
  review → rejected   when refused (human); failed → review for retry
  running → transferred  when work crosses a project boundary:
             CloudEvent card → target agent's inbox → its own queued state
```

Rules:
- No state transition without a sealed journal row (F0).
- A debrief (`.agent_learning` record) is mandatory between `completed` and `closed`.
- Fresh context per state: rehydrate from the journal, never carry a long window.

## Review policy

Order: **deny > ask > allow.** Check, in order:

1. Scope: the agent only touches its own project; crossing projects requires a transfer.
2. Pricing: expensive (pro) provider only off-peak; at peak → `deferred` or flash fallback.
3. Egress: filter text before any paid API; refused trees never cross.
4. Irreversibility: push, email, tunnel DNS, data drop, SMS → human confirmation first.
5. Blast radius: Critical/High findings, shared components, money/safety changes → human review.

## Transfers (any agent → any agent)

Control-mediated, file-only. Source agent writes a CloudEvent-shaped card
(`specversion, type, source, id, data`) into the target's `inbox/`; the target's watcher
(one-shot, never polling) wakes it; it hydrates its own taxonomy + the card, works, returns
a result card, debriefs. Control seals `handoff.initiated` / `handoff.completed`. Agents
never write outside their own project except through control. Lineage (`parent:` in the
manifest) is capped and recorded.

## Child-agent forking (human-gated)

- Human-gated: an agent *proposes* a child; a human approves it. Children are born untrusted.
- Narrower scope, never wider: a child inherits F0–F4 with narrower scope — smaller model,
  cheaper run.
- Expiry + re-validation (a scoped skill goes stale when its library changes).
- Per-lineage rollback; population cap (2 children per parent); kill-switch freeze.
- Elevation (untrusted → … → saint) is earned on sealed debriefs + monitor verdicts.

## Skill sizing (deterministic, zero tokens)

`manifest tier (local|flash|pro)` → device capacity → local-model presence →
`route(kind)` (local-first; flash after egress; pro only off-peak) → context caps
(file / dir / root) → egress filter.

## Directives (control layer; ship keeps D0–D7)

| ID | Directive |
|---|---|
| C0 | Control ledger chain intact; every dispatch/transfer/completion sealed (F0) |
| C1 | No workflow starts a provider before `workflow.cleared` |
| C2 | No cross-project writes outside the transfer protocol |
| C3 | No expensive provider at peak; pro only off-peak |
| C4 | No `closed` without a debrief record |
| C5 | Child-agent forking is human-gated (propose → approve) |
| C6 | All code shipping goes through frontier plan → apply (never bypass) |
| C7 | No xAI/Anthropic/OpenAI API calls anywhere in the platform |
| C8 | All cloud calls are api.deepseek.com; pro off-peak only |
| C9 | No new runs when session context/cache is ≥95% — soft-refuse, seal a handoff, open a new session |

## Session budget (C9)

When a session's context or cache reaches **95%**, soft-refuse new runs: stop starting work, write handoff artifacts (`HANDOFF.md` + a learning debrief), and continue in a new session. Deterministic check: `session_budget(usage_pct)` in the Python reference, Go `internal/sessionbudget`, and review before dispatch.

## Providers

`local` (on-host model, 0 cloud tokens) · `deepseek-flash` (anytime, has Vision) ·
`deepseek-v4-pro` (off-peak only, no Vision) · `dsh` (harness sessions). One workflow, any
provider — the agent's knowledge lives in files, never in a session.

**Ship safely. Spend less.**

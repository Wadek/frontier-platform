# Control Tower — doctrine (frontier-control · Milkcow)

The tower is the air-traffic control of the fleet. Pilots (agents) fly planes (project
directories); the tower routes, clears, hands off, lands, and debriefs. It never does the
work itself, and a model never decides which model to use.

## Components

| Component | Job |
|---|---|
| taxonomist | manage per-plane taxonomy branches (`tower taxonomy …`) |
| data index | registry of planes + pilots (`tower roster`), derived from the file tree |
| weather | peak/off-peak gating + daily pricing-page check (`tower weather`) |
| job service | scheduled one-shots via Task Scheduler — never polling (`tower ops jobs`) |
| workflow engine | FSM over an append-only journal; every state rehydrates fresh context |
| review | the conscience gate: scope, weather, egress, deny>ask>allow (`tower workflow review`) |
| runways | thin harness adapters: same workflow, `--runway waka-cli\|deepseek-flash\|deepseek-v4-pro\|dsh` |
| monitor | replay the ledger against directives (`tower monitor`) |
| prove | Proven/Theory harness — no claim without a passing test (`tower prove`) |

## Workflow FSM

```
queued → review → taxiing (hydrate own taxonomy) → airborne → landed → debriefed → filed
  review → holding   when pro runway requested at peak (job service reschedules off-peak)
  review → rejected  when refused (human); failed → review for retry
  airborne → handed_off  when work crosses a plane boundary:
             CloudEvent card → target pilot's inbox → its own queued state
```

Rules:
- No state transition without a sealed journal row (F0).
- A debrief (`.agent_learning` record) is mandatory between `landed` and `filed`.
- Fresh context per state: rehydrate from the journal, never carry a long window.

## Clearance (review)

Order: **deny > ask > allow.** Check, in order:

1. Scope: the pilot only touches its own plane; crossing planes requires a handoff.
2. Weather: pro runway only off-peak; at peak → `holding` or flash fallback.
3. Egress: `filter_text` before any paid API; refused trees never cross (see waka-cli `aware/egress.py`).
4. Irreversibility: push, email, tunnel DNS, data drop, SMS → human confirmation first.
5. Blast radius: Critical/High findings, shared components, money/safety changes → human review.

## Handoffs (any pilot → any pilot)

Tower-mediated, file-only. Source pilot writes a CloudEvent-shaped card
(`specversion, type, source, id, data`) into the target's `inbox\`; the target's watcher
(one-shot, never polling) wakes it; it hydrates its own taxonomy + the card, works, returns
a result card, debriefs. The tower seals `handoff.initiated` / `handoff.completed`. Pilots
never write outside their own plane except through the tower. Lineage (`parent:` in the
manifest) is capped and recorded.

## Speciation (from the SSC kit, `original_control_plane_ideas`)

- Human-gated: a pilot *proposes* a child; a human approves it. Children are born untrusted.
- Narrower scope, never wider: a child inherits F0–F4 with narrower scope — smaller model,
  cheaper run.
- Expiry + re-validation (a scoped skill goes stale when its library changes).
- Per-lineage rollback; population cap (2 children per parent); kill-switch freeze.
- Elevation (untrusted → … → saint) is earned on sealed debriefs + monitor verdicts.

## Skill sizing (deterministic, zero tokens)

`manifest tier (local|flash|pro)` → device VRAM → local-model presence (roster) →
`route(kind)` (local-first; flash after egress; pro only off-peak) → context caps
(file 2 KB / dir 3 KB / root 4 KB) → egress filter. Reuse waka-cli `aware/*` wholesale.

## Directives (tower layer; ship keeps D0–D7)

| ID | Directive |
|---|---|
| T0 | Tower ledger chain intact; every dispatch/handoff/landing sealed (F0) |
| T1 | No workflow starts a runway before `workflow.cleared` |
| T2 | No cross-plane writes outside the handoff protocol |
| T3 | No expensive runway at peak; pro only off-peak |
| T4 | No `filed` without a debrief record |
| T5 | Speciation is human-gated (propose → approve) |
| T6 | All code shipping goes through frontier plan → apply (never bypass) |
| T7 | No xAI/Anthropic/OpenAI API calls anywhere in the fleet |
| T8 | All cloud calls are api.deepseek.com; pro off-peak only |
| T9 | No new runs when the session context/cache is ≥95% full — soft-refuse, seal the handoff, open a new session |

## Session budget (the 95% rule — T9)

When a session's context or cache reaches **95%**, the agent **soft-refuses** new runs:
it stops starting work, **updates the handoff** (`HANDOFF.md` for the plane + a
`.agent_learning` debrief record), and **insists on opening a new session** to continue.
Soft = the handoff artifacts are required before stopping; never a silent stop, never a
hard crash into a full cache. The decision is deterministic:
`session_budget(usage_pct)` (Python `reference/tower_ref.py` · Go
`internal/sessionbudget` · Haskell `Tower.SessionBudget`) — at ≥95% it returns
`refuse_new_runs + require_handoff + require_new_session`. The review gate consults it
before every dispatch, and every pilot obeys it between turns.

## Runways

`waka-cli` (local Ollama, 0 tokens) · `deepseek-flash` (anytime, has Vision) ·
`deepseek-v4-pro` (off-peak only, no Vision) · `dsh` (this harness). One workflow, any
runway — the pilot's knowledge lives in files, never in a session.

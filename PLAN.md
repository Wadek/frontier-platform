# frontier-fleet — bootstrap plan (condensed from planning session 2026-09-14)

Read this first at execution time. Session debrief: `sessions\.agent_learning.jsonl`.
Do NOT re-derive anything in this file. Zero model spend to reach here — everything below
is already decided and verified.

## Mission

Cost-efficient, safe development with frontier models. Umbrella for the control plane:
**frontier-ship** (the gate) + **frontier-control** (the tower) + **frontier-insight** (AI-usage
observatory) + **frontier-ops** (docker/uptime/IP ground crew). "The fleet ships safely and
flies cheap."

## Locked decisions

- Umbrella: **frontier-fleet** · control codename **Milkcow** (Type XIV supply U-boat).
- Pilots at wakalabs root: `.agent_<plane>\` — taxonomy branch per plane:
  `taxonomy.yaml` (manifest: version, created_by, domain, runways, lineage, tier, expiry) +
  `knowledge\` (architecture.md, mappings.md, brief.md — each pins repo/commit/patterns) +
  `compositional_skills\` (grounded\ + ungrounded, SKILL.md bodies, qna.yaml-style frontmatter) +
  `sessions\.agent_learning.jsonl` (append-only; records: session, plane, pilot, runway, ts,
  tldr, learned[], skills_proposed[], mappings_updated[], next[]) + `inbox\` (CloudEvent-shaped
  handoff cards {specversion,type,source,id,data}).
- **One repo**: `Wadek/frontier-fleet` monorepo — `ship\` (history imported via `git subtree`;
  old `frontier-ship` remote archived after human confirm), `control\` (binary `tower`),
  `insight\`, `ops\`, `runtime\` (hooks, bin junction, ledgers, PIPELINE.md), and all
  `.agent_<plane>\` pilot packs (fleet assets). App products (food, opengym, satokori, …)
  keep their own repos, remotes, CI, and branch maps. `D:\frontier` junction re-points to
  `frontier-fleet\runtime`; global `core.hooksPath` updated. deepseek_harness, immich,
  home-assistant, adguard, ollama-data are external/vendor — never merged.
- v1 = CLI + files only (no web cockpit).
- Red Hat vocabulary adopted (see TERMS.md to be written): taxonomy, knowledge vs
  compositional skills, grounded/ungrounded, seed_examples, teacher/student + data generate,
  orchestrator, workflow + workflow instances, data index, job service, review, CloudEvents.
  ilab verbs → `tower taxonomy init|diff|add|approve|retire`, `tower data generate --plane <p>`.
- **Single-vendor model policy**: cloud = api.deepseek.com only (deepseek-flash anytime,
  supports Vision; deepseek-v4-pro off-peak only, no Vision). Local = Ollama on RTX 2080 8 GB.
  xAI/Grok and Anthropic/Claude removed everywhere. T-directives: T7 no xAI/Anthropic/OpenAI
  calls; T8 all cloud = api.deepseek.com, pro off-peak only.
- **Weather** (official pricing page, fetched 2026-09-14): peak = Mon–Fri 01:00–04:00 and
  06:00–10:00 UTC; everything else off-peak (50% price). Daily one-shot check of the pricing
  page (DeepSeek says to check regularly; V4 Pro service continues past 2026-09-14, unchanged
  billing). Pro runway at peak → holding state or flash fallback; off-peak queues via
  Task Scheduler (job service), never polling.
- **Deprecate**: caravans, human-atlas, pigeon, phylax, watermarks-remover →
  `tower ops decommission <plane>`: salvage (harvest skills/learnings into `.agent_<plane>`)
  → stop/remove containers → remove gateway/tunnel routes → back up data (**phylax vault DB
  requires explicit human confirm**) → mark plane retired in registry. Salvage first:
  caravans' token_tracker (→ fleet standard), watermarks' install_skill.py (→ runway-adapter
  reference) and core hygiene API (→ re-homed to `frontier-fleet\runtime\hygiene\`; ML
  harnesses retire with the project).

## Tower command surface

```
tower taxonomy init|diff|add|approve|retire   tower data generate --plane <p> (teacher pass, off-peak)
tower workflow run|status|review|list         tower roster (data index)
tower ai audit|improve <plane>                tower ops health|docker|ip|decommission|jobs
tower runways list|test                       tower weather show|check
tower monitor (T-directives replay)           tower prove (Proven/Theory harness, from SSC emulate.py)
```

Components: taxonomist · data index · weather · job service · workflow engine (FSM over
append-only journal; fresh context per state) · review (clearance/conscience: scope, weather,
egress, deny>ask>allow, Critical/High/irreversible → human) · runway adapters · monitor · prove.
Shared ledger = append-only JSONL (F0), frontier-compatible; ship gate rows + tower rows.

## Workflow FSM

`queued → review → taxiing (hydrate own taxonomy) → airborne → landed → debriefed → filed`
+ `holding` (pro@peak, job service reschedules) + `handed_off` (crosses plane boundary:
CloudEvent card → target pilot's inbox → its own queued state) + `failed` (retry via review)
+ `rejected` → end. Debrief mandatory before `filed`.

## Skill sizing (deterministic, zero tokens — reuse waka-cli aware/*)

pilot manifest tier (local|flash|pro) → `devices.py` VRAM inventory → `roster.py` local-model
presence → `matrix.route(kind)` (local-first; flash after egress; pro only off-peak) →
`hydrate.py` context caps (file 2 KB / dir 3 KB / root 4 KB) → `egress.filter_text` before any
paid API (REFUSE list + aliases). A model never decides which model to use.

## Speciation (from original_control_plane_ideas SSC kit)

Human-gated propose → approve. Children born untrusted, narrower scope inherits F0–F4
(never wider), expiry + re-validation, per-lineage rollback, population cap (2/parent),
kill-switch freeze. Elevation earned on sealed debriefs + monitor verdicts
(untrusted → … → saint).

## Fleet AI-usage survey (done, read-only; source for Phase-2 census + improvement queue)

- food: local-only culina:7b recipes; prices never from model. Adds: accepted/rejected flag,
  deterministic fallback, agent-card (copy satokori's).
- satokori: Claude CLI tips (10 isolated calls; likely broken in Docker) → replace with local
  7B draft + deepseek-flash final; remove ANTHROPIC_API_KEY; keep ai-plugin.json/MCP/agent-card.
- opengym: coach.js migrates xAI grok-4.6 → deepseek-flash (pool logic kept, price constants
  swap, `DEEPSEEK_COACH_BUDGET_USD`, delete XAI_* env). Implement **deepseek-verify** (pro,
  off-peak, sampled) — the missing "double-check first 10". Single-source the coach doctrine
  (one SKILL.md). Record DeepSeek usage into coach-usage.json.
- polly-engine: CPU CV pipeline stays; optional grok vision → deepseek-flash Vision + spend ledger.
- waka-net: infra plane (Ollama + waka-bot UI + vigil). Wire vigil token dashboard to tower
  data index (metric: local:deepseek ratio); dedupe waka-bot system prompt.
- No-AI: tasks, books, pigeon, phylax, home-assistant, adguard, human-atlas (audit = one-liner).
- Vendors removed fleet-wide: xAI, Anthropic. Ollama + DeepSeek only.

## Phases (cost-gated)

- **Phase 0** (on go, cheap): debrief record already written (this session).
- **Phase 1** (cheap/local runways): restructure tree (ship move, junction, hooksPath,
  reference updates in CONTEXT.md, ai_learning 02/09/10/12, MANIFEST.json, PIPELINE.md,
  prompts/local-system.md ×2); scaffold `control\` (English doctrine + TERMS.md +
  ARCHITECTURE_ASCII.md, Haskell FSM/clearance witnesses, Go runtime; F5 consilience law);
  runway adapters (waka-cli, deepseek, dsh); weather + daily check; `.agent_learning` spec +
  templates; watcher scripts. Ship everything via frontier plan→apply→push (dogfood), tests
  same turn.
- **Phase 2** (pro/max, off-peak only): `tower census` over all planes → `.agent_<plane>\`
  packs + `docs/AI_USAGE.md` per project. One plane per run, journal-driven, resumable,
  ledger-sealed. Fits the ~39 h off-peak window opening 2026-09-14 10:00 UTC.
- **Phase 3**: dogfood — one real task end-to-end incl. a cross-agent handoff; monitor clean.
- **Phase 4**: speciation/elevation loop, monitor-watch integration, `ai_learning/13-fleet.md`
  + CONTEXT.md updates, release.

## Habitat rules (hold throughout)

Feature branch, never main · frontier plan/apply, never --no-verify/FRONTIER_SOFT · no
polling (Task Scheduler + one-shot watchers) · tests verbose same turn · confirm before push,
email, DNS, or dropping data · local-first · secrets never in agent packs, chat, or git.

## First steps at go

1. Read this file + `sessions\.agent_learning.jsonl`.
2. Restructure: `git init` at `frontier-fleet\`; subtree-import ship history into `ship\`;
   re-point `D:\frontier` junction; update `core.hooksPath`; archive old remote after confirm.
3. Scaffold `control\` in the same repo (feature branch, English/Haskell/Go skeleton + tests)
   → plan/apply → push → PR (human merges).
4. `weather.yaml` with the exact peak window + scheduled daily pricing check.
5. Queue the five decommissions via `tower ops` (phylax data backup gated on human confirm).
6. Off-peak: run `tower data generate` (teacher pass) → `tower census` per plane.

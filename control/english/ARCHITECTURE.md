# Architecture

Four views. PlantUML sources: `docs/architecture/*.puml`.

## 1. System context

```
                        ┌──────────────────────────────────────────┐
                        │                  HUMAN                   │
                        │  control console · review approvals ·    │
                        │  child-agent approvals · git merges      │
                        └──────────────────────┬───────────────────┘
                         orders / reviews     │ approve
┌──────────────────────────────────────────────▼────────────────────────────────────────────┐
│                              FRONTIER-PLATFORM                                            │
│                                                                                            │
│  ┌── frontier-control (binary: control) ───────────────────────────────────────────────┐  │
│  │  ┌────────────┐ ┌─────────────┐ ┌────────────┐ ┌─────────────┐ ┌────────────────┐   │  │
│  │  │ taxonomist │ │  data index │ │  pricing   │ │ job service │ │ workflow engine│   │  │
│  │  └─────┬──────┘ └──────┬──────┘ └─────┬──────┘ └──────┬──────┘ └───────┬────────┘   │  │
│  │  ┌─────▼──────┐ ┌──────▼──────┐ ┌─────▼──────┐ ┌──────▼──────┐ ┌───────▼────────┐   │  │
│  │  │ teacher    │ │  review     │ │ providers  │ │  monitor    │ │    prove       │   │  │
│  │  │ data gen   │ │ (clearance) │ │  adapters  │ │ directives  │ │ Proven/Theory  │   │  │
│  │  └────────────┘ └─────────────┘ └────────────┘ └─────────────┘ └────────────────┘   │  │
│  └──────────────────────────────────┬──────────────────────────────────────────────────┘  │
│  ┌───────────────────┐      ┌───────▼────────┐      ┌──────────────────────────────────┐ │
│  │ frontier-platform │◄─────┤ shared ledger  │◄─────┤ runtime: hooks · bin · PIPELINE  │ │
│  │ plan→apply→       │ gate │ append-only    │      │                                  │ │
│  │ push · G/H/R/O    │ only │ JSONL · F0     │      │                                  │ │
│  └───────────────────┘      └────────────────┘      └──────────────────────────────────┘ │
└──────────────────────────────────────────────┬────────────────────────────────────────────┘
                                               │ workflow runs · transfer cards
        ┌──────────────────────────────────────┼───────────────────────────────────────┐
        │              PROJECTS — top-level working directories                         │
        │   each agented by .agent_<project>/ at the workspace root                      │
        └──────────────────────────────────────┼───────────────────────────────────────┘
                                               │ dispatch · skill sizing
┌──────────────────────────────────────────────▼───────────────────────────────────────┐
│   PROVIDERS (model gateway) + LOCAL AWARENESS                                          │
│   local → on-host model (0 cloud tokens)                                               │
│   deepseek-flash (anytime, Vision) · deepseek-v4-pro (off-peak only)                   │
│   sizing: capacity · presence · route(kind) · context caps · egress · is_peak()        │
└───────────────────────────────────────────────────────────────────────────────────────┘
```

## 2. Taxonomy tree (one branch per project)

```
<workspace>/
├── frontier-platform/       ship/ · control/ · runtime/ · .agent_frontier-platform/
├── .agent_<project>/        taxonomy.yaml · knowledge/ · compositional_skills/
│                            (grounded/ + ungrounded) · sessions/.agent_learning.jsonl · inbox/
│                            compositional_skills/learning employed by other skills
└── .agent_<other>/ …
```

## 3. Workflow lifecycle (FSM)

```
(new task) → queued → review → preparing → running → completed → debriefed → closed → [*]
              review → deferred (pro@peak; job service) → review (off-peak)
              review → rejected → [*]         failed → review
              running → transferred (CloudEvent card → target inbox → its queued)
```

## 4. Transfer (sequence)

```
Human → control: workflow run
control → agent A: dispatch + review (provider, pricing)
A: hydrate own taxonomy · work · hits project boundary
A → control: transfer request
control → agent B: CloudEvent card → .agent_<B>/inbox/
control → ledger: seal handoff.initiated
watcher (one-shot) → B: wake · hydrate own taxonomy + card · work
B → control: result card · write .agent_learning.jsonl
control → ledger: seal handoff.completed + debrief
control → Human: review (Critical / irreversible / child-agent forking)
```

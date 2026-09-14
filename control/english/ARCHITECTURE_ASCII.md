# Architecture (ASCII)

Four views. PlantUML sources: `docs/architecture/*.puml`.

## 1. System context

```
                        ┌──────────────────────────────────────────┐
                        │             HUMAN (Wade)                 │
                        │  tower console · review approvals ·      │
                        │  speciation approvals · git merges       │
                        └──────────────────────┬───────────────────┘
                         orders / reviews     │ approve
┌──────────────────────────────────────────────▼────────────────────────────────────────────┐
│                FRONTIER-FLEET · D:\wakalabs\frontier-fleet\                               │
│                                                                                            │
│  ┌── frontier-control (binary: tower · codename Milkcow) ──────────────────────────────┐  │
│  │  ┌────────────┐ ┌─────────────┐ ┌────────────┐ ┌─────────────┐ ┌────────────────┐   │  │
│  │  │ taxonomist │ │ data index  │ │  weather   │ │ job service │ │ workflow engine│   │  │
│  │  └─────┬──────┘ └──────┬──────┘ └─────┬──────┘ └──────┬──────┘ └───────┬────────┘   │  │
│  │  ┌─────▼──────┐ ┌──────▼──────┐ ┌─────▼──────┐ ┌──────▼──────┐ ┌───────▼────────┐   │  │
│  │  │ teacher    │ │  review     │ │  runways   │ │  monitor    │ │    prove       │   │  │
│  │  │ data gen   │ │ (clearance) │ │  adapters  │ │ T-directives│ │ Proven/Theory  │   │  │
│  │  └────────────┘ └─────────────┘ └────────────┘ └─────────────┘ └────────────────┘   │  │
│  └──────────────────────────────┬──────────────────────────────────────────────────────┘  │
│  ┌───────────────┐      ┌───────▼────────┐      ┌──────────────────────────────────┐     │
│  │ frontier-ship │◄─────┤ shared ledger  │◄─────┤ runtime: hooks · bin · PIPELINE  │     │
│  │ plan→apply→   │ gate │ append-only    │      │ (D:\frontier junction re-points  │     │
│  │ push · G/H/R/O│ only │ JSONL · F0     │      │  here; old paths keep working)   │     │
│  └───────────────┘      └────────────────┘      └──────────────────────────────────┘     │
└──────────────────────────────────────────────┬────────────────────────────────────────────┘
                                               │ workflow runs · handoff cards
        ┌──────────────────────────────────────┼───────────────────────────────────────┐
        │              PLANES — top-level wakalabs dirs                               │
        │   food · opengym · satokori · tasks · polly-engine · waka-net · …            │
        │   each piloted by .agent_<plane>\ at D:\wakalabs root                        │
        └──────────────────────────────────────┼───────────────────────────────────────┘
                                               │ dispatch · skill sizing
┌──────────────────────────────────────────────▼───────────────────────────────────────┐
│   RUNWAYS (model gateway) + LOCAL AWARENESS (waka-cli aware/*)                        │
│   waka-cli → Ollama qwen2.5-coder:64k on RTX 2080 8 GB (0 tokens)                     │
│   deepseek-flash (anytime, Vision) · deepseek-v4-pro (off-peak only)                  │
│   sizing: devices.py VRAM · roster.py presence · matrix.route(kind)                   │
│            hydrate.py caps (2/3/4 KB) · egress.filter_text · is_peak()                │
└───────────────────────────────────────────────────────────────────────────────────────┘
```

## 2. Taxonomy tree (one branch per plane)

```
D:\wakalabs\
├── frontier-fleet\          ship\ · control\ · runtime\ · .agent_frontier-fleet\
├── .agent_food\             taxonomy.yaml · knowledge\ · compositional_skills\
│                            (grounded\ + ungrounded) · sessions\.agent_learning.jsonl · inbox\
└── .agent_opengym\ … (~20 pilots)
```

## 3. Workflow lifecycle (FSM)

```
(new task) → queued → review → taxiing → airborne → landed → debriefed → filed → [*]
              review → holding (pro@peak; job service) → review (off-peak)
              review → rejected → [*]         failed → review
              airborne → handed_off (CloudEvent card → target inbox → its queued)
```

## 4. Handoff (sequence)

```
Human → tower: workflow run
tower → pilot A: dispatch + review (runway, weather)
A: hydrate own taxonomy · work · hits plane boundary
A → tower: handoff request
tower → pilot B: CloudEvent card → .agent_<B>\inbox\
tower → ledger: seal handoff.initiated
watcher (one-shot) → B: wake · hydrate own taxonomy + card · work
B → tower: result card · write .agent_learning.jsonl
tower → ledger: seal handoff.completed + debrief
tower → Human: review (Critical / irreversible / speciation)
```

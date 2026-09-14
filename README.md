# frontier-fleet

The wakalabs control plane. One monorepo, four capabilities:

| Dir | Capability | Binary / job |
|---|---|---|
| `ship/` | **frontier-ship** — the gate: plan → apply → push, hash-chained ledger (F0), monitor D0–D7 | `frontier` |
| `control/` | **frontier-control** (codename Milkcow) — the tower: taxonomy, workflows, runways, weather, speciation | `tower` |
| `control/` (ai family) | **frontier-insight** — AI-usage observatory: audit + improve | `tower ai` |
| `control/` (ops family) | **frontier-ops** — ground crew: docker, health, IP drift, decommission | `tower ops` |

Mission: **cost-efficient, safe development with frontier models.** *"The fleet ships safely and flies cheap."*

- `PIPELINE.md` — the permanent ship process (canonical spec)
- `PLAN.md` — bootstrap plan this repo was born from
- `sessions/.agent_learning.jsonl` — append-only session debriefs (dogfood)
- `runtime/` — live gate: hooks, binaries, ledgers (git-ignored; `D:\frontier` junctions here)

## Laws

The ship's axioms F0–F5 bind everything here: evidence before remote effect, never main,
never bypass the gate, English/Haskell/Go consilience. The tower adds directives T0–T8
(`control/english/CONTROL_TOWER.md`).

## Build & test

- Ship: Go toolchain (see `ship/README.md`).
- Tower: Go runtime sources in `control/`; a Python-stdlib reference with tests in
  `control/reference/` (`python -m unittest discover -s control/tests -v`) runs on any host.
- Haskell witnesses: `control/haskell/` (proof surface; build where GHC exists).

## Gate (this repo dogfoods itself)

`lefthook.yml` + `scripts/frontier-gate.ps1` wire the global `core.hooksPath` dispatchers:
feature branch → `frontier plan` → `frontier apply` → push. Never `--no-verify`, never main.

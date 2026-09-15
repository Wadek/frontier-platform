# frontier-platform

Local-first tooling to **ship safely** and **spend less** on model usage.

One monorepo, four capabilities:

| Dir | Capability | Binary / job |
|---|---|---|
| `ship/` | **frontier-ship** — gate: plan → apply → push, hash-chained ledger (F0), monitor D0–D7 | `frontier` |
| `control/` | **frontier-control** — taxonomy, workflows, providers, pricing, child-agent forking | `control` |
| `control/` (ai) | **frontier-insight** — AI-usage audit + improve | `control ai` |
| `control/` (ops) | **frontier-ops** — docker, health, IP drift, decommission | `control ops` |

- [`PIPELINE.md`](PIPELINE.md) — permanent ship process
- [`control/`](control/) — coordinator doctrine and reference implementation
- [`hygiene/`](hygiene/) — optional AI-provenance inspect service for `frontier hygiene`
- `runtime/` — live gate binaries, hooks, ledgers (gitignored; point `FRONTIER_RUNTIME` here)

## Laws

Ship axioms F0–F5 bind shipping: evidence before remote effect, never push `main`, never bypass the gate, English/Haskell/Go consilience. Control adds directives C0–C9 (`control/english/CONTROL.md`).

## Build and test

- Ship: Go toolchain (see `ship/README.md`).
- Control: Go under `control/`; Python stdlib reference + tests:

```powershell
python -m unittest discover -s control/tests -v
```

- Haskell witnesses: `control/haskell/` (where GHC exists).

## Gate (this repo dogfoods itself)

`lefthook.yml` + `scripts/frontier-gate.ps1` dispatch: feature branch → `frontier plan` → `frontier apply` → push. Never `--no-verify`, never `main`.

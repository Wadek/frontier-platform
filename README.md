# frontier-platform

Local-first **security (S)**, **AI usage (L4)**, and **delivery (L7)** controls so AI-assisted code can reach production without skipping checks.

Verbose habitat plan: `D:\wakalabs\docs\CODE_TO_PRODUCTION_PLAN.md` (diagrams in `docs/architecture/`).

## Branches (apps)

`feat/*` → `dev` (local serve) → `alpha` (versioned) → `main` (human merge, versioned prod).

## Ship (`ship/`)

```powershell
frontier plan          # or: frontier plan --json
frontier apply
frontier release-check # prod deploy authorize (not a push)
git push               # hooks re-run plan/apply
```

Never push `main` via Frontier. Prod deploy uses `release-check`, then app `deploy.ps1`.

Release templates: `ship/templates/release/`. Process: `ship/english/RELEASE_PROCESS.md`.  
Token thrift: printed on plan/apply OK; skill `grok/token-thrift`.

## Control (`control/`)

First-class AI usage plane: `control pricing|roster|check`. Other verbs reserved (not all implemented in v0.1).

## Standard vs custom

| standard | custom |
|----------|--------|
| git, GitHub PR/Actions/Environments, Docker/Podman, Gitleaks, Trivy, Checkov, Semgrep (CI), Playwright, ZAP | Frontier plan/apply/release-check/ledger/hooks, deploy templates, control CLI, token thrift |

## Build

- Ship: Go (`ship/`). Place `frontier.exe` under `runtime/bin` (or `FRONTIER_RUNTIME`).
- Control: `python -m unittest discover -s control/tests -v`

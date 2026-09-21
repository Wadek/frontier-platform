# Release process — feat → dev → main

Standard merge and deploy. Open-source tools do the scans. Frontier authorizes the laptop push and the production deploy. GitHub branch protection owns “no direct push to `main`.”

## Branches

```text
feat/<slug> → PR → dev → PR → main
```

| Branch | Serves | Who merges | Version |
|--------|--------|------------|---------|
| feat/* | laptop | author (PR into `dev`) | — |
| dev | integration | after Stages 1–2 green | — |
| main | production | **human** | `prod-YYYYMMDD-n` via deploy |

There is no alpha channel in this process.

## Who does what

| Concern | Owner |
|---------|--------|
| Require PR + checks on `main` | **GitHub** branch protection |
| Exam + ledger before a laptop push | Frontier `plan` → `apply` |
| Secrets / SAST / SCA / IaC / tests | GitHub Actions (`verify.yml`) using OSS CLIs |
| Mutate production | App `deploy.ps1` (or equivalent) |
| “This commit may go to production” | Frontier `release-check` inside that deploy |

Laptop deny-commit-on-main is a hook catch, not a separate subsystem.

## Stages (no LLM)

| Stage | When | Tools |
|-------|------|--------|
| **0** | `git push` of `feat/*` or `dev` | Frontier `plan` → `apply` (OWASP v0 High/Critical, dirty tree, ledger) |
| **1** | PR into `dev` or `main` | Gitleaks, Semgrep, Trivy fs, Checkov, ESLint/Ruff |
| **2** | same PRs | Stack tests (`pytest -v`, `node --test`, `npm test`, `go test`) |
| **3** | PR | container build, Trivy image, Checkov on Dockerfile/compose |
| **4** | optional later | Playwright, ZAP, light k6 |
| **5** | after human merge to `main` | `deploy.ps1`: backup, up, health, rollback. `frontier release-check` first. |

Copy `ship/templates/release/` into the app on **first deploy**. `frontier ready` lists what is still missing.

## First deploy (detect, do not mutate)

From the app tree:

```text
frontier ready
```

JSON + text: pass / fail / advise per item, with a fix hint. Frontier does not write the app. An agent or human applies the hints, then re-runs `ready`.

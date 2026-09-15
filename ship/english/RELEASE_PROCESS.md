# Release process — feat → dev → alpha → main

Habitat path: see `D:\wakalabs\docs\CODE_TO_PRODUCTION_PLAN.md` (verbose). This file is the platform short form.

## Branches

```text
feat/<slug> → dev → alpha → main
```

| Branch | Serves | Merge | Version |
|--------|--------|-------|---------|
| feat/* | laptop | author | — |
| dev | local dev server | AI after Stage 1–2 | — |
| alpha | staging URL | AI after Stage 1–4 | `alpha-YYYYMMDD-n` |
| main | production | **human** | `prod-YYYYMMDD-n` |

## Frontier (custom) inside standard CI

- Every laptop push: `frontier plan` → `apply` (hooks).
- Prod deploy: `frontier release-check` (tolerates only `refuse_push_main`).
- `frontier plan --json` / `release-check --json` for automation.

## Stages (standard tools)

1. Static: Gitleaks, Semgrep, Trivy fs, lint + Frontier Guard  
2. Unit: project test runners + build  
3. Integration: compose/CI + Trivy image + Checkov  
4. Alpha: deploy + Playwright + ZAP + light k6  
5. Prod: release-check + deploy + smoke + rollback  

## Templates

Copy `ship/templates/release/` into an app (wakagym first). Set `RUNTIME=docker` or `podman`.

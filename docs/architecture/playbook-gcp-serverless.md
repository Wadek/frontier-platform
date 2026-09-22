# GCP Serverless Playbook (FastAPI + Cloud Run)

Reusable patterns for shipping a FastAPI app on Google Cloud Run with Firebase Hosting as edge ingress. No project-specific names.

## Secrets

- Never mount Secret Manager volumes onto the app working directory (`/app`). That shadows the tree and wipes code at boot.
- Mount to `/secrets/` (e.g. `/secrets/.env`) and symlink into the app path in the image build: `ln -s /secrets/.env /app/.env`.

## Configuration

- Cloud Run injects system env vars (`PORT`, `K_SERVICE`, …). Pydantic settings must use `SettingsConfigDict(extra="ignore")` (or equivalent) so those vars do not crash boot.
- Keep app secrets in Secret Manager / env; do not commit API keys or tokens into source.

## Containers

- Prefer array-form `CMD` / `ENTRYPOINT` (e.g. `CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]`).
- Avoid production reliance on `entrypoint.sh` when Windows CRLF or missing shebang can yield opaque `no such file or directory` failures.
- Run DB migrations on boot only if idempotent and fast enough for cold start (e.g. `alembic upgrade head` before uvicorn).

## Background work

- Do not use FastAPI `BackgroundTasks` for work that must finish after the HTTP response on Cloud Run: CPU can drop to near zero when the request completes.
- Use a task queue (Cloud Tasks) calling an authenticated worker HTTP endpoint instead.

## Networking

- Native Cloud Run custom domain mappings are unavailable or awkward in some regions.
- Prefer Firebase Hosting (or another edge proxy) with rewrite rules to the Cloud Run service for HTTPS and custom domains.

## CI / deploy

- Prefer Workload Identity Federation (OIDC) from GitHub Actions over static JSON keys.
- When GitHub-hosted minutes are exhausted, keep the same workflows and switch `runs-on` to self-hosted labels; do not leave `ubuntu-latest` jobs burning the quota.
- On Windows self-hosted runners under LocalSystem, use `shell: powershell` (not `pwsh` unless PowerShell 7 is on the machine PATH).

## API UX

- For endpoints that trigger irreversible external actions (SMS, email, payments), support `?dry_run=true` (or body flag) so clients can preview before commit.

## Scanners

- Gitleaks = secrets; Trivy fs = `--scanners vuln` when gitleaks already covers secrets; Checkov = CLI (`checkov` / `checkov.cmd`), never `python -m checkov`.
- Emit JSON reports and open deduped GitHub issues from the runner; local-first agents may propose hotfixes under human control.

# Docker Compose GitHub Actions runner

Linux self-hosted runner with **pre-baked** Stage 1 scanners (gitleaks, trivy, checkov, semgrep). Prefer this over Windows service runners and over `pip install` every job.

## Why Docker Desktop

- No host/WSL sudo for day-to-day operation.
- Avoids Windows LocalSystem `pwsh` / `checkov.cmd` / cp1252 Semgrep JSON friction.
- Matches `runs-on: [self-hosted, Linux, X64]` + `shell: bash` in `verify.yml.example`.

## Setup

1. Copy this directory into the app as `ops/github-runner/` (or keep a shared habitat copy).
2. Copy `.env.example` → `.env` and set `REPO_URL` + `RUNNER_TOKEN`.
3. `docker compose up -d --build`
4. Confirm the runner is **Idle/Online** in the repo Actions settings.
5. Point workflows at `[self-hosted, Linux, X64]`.

## Semgrep note

Pin `setuptools==75.8.2` (or another release that still provides `pkg_resources`) inside the semgrep venv. Newer setuptools alone can break `semgrep --version`.

## Checkov note

Invoke the console script `checkov` on PATH. Never `python -m checkov` (no `__main__`).

## CKV_DOCKER_8

This image intentionally ends as root so the upstream runner entrypoint can register and (optionally) use the Docker socket. The Dockerfile documents a Checkov skip with that reason. App images should still use non-root `USER` + `HEALTHCHECK`.

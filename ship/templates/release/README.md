# Release templates

Generic verify / deploy / rollback for apps that use Frontier.

## Copy into the app

1. `verify.yml.example` → `.github/workflows/verify.yml` (hard-fail Stage 1 + issue reporter + Stage 2 tests).
2. `scripts/ci/report_scanner_issues.py` → `scripts/ci/report_scanner_issues.py`.
3. `pull_request_template.md` → `.github/pull_request_template.md` (token report + Closes #n).
4. Optional GCP path: `cloud-deploy.yml.example` → `.github/workflows/cloud-deploy.yml` (fill WIF / project / region).
5. `cleanup-merged.yml.example` → `.github/workflows/cleanup-merged.yml` (close/delete issues + delete head branch after merge to `dev`).
6. Optional Compose deploy scripts: `deploy.yml.example`, `runtime.ps1`, `env-table.example.ps1`.
7. Self-hosted Linux runner: copy `../runner-docker/` to `ops/github-runner/` and `docker compose up -d --build`.

## Stage 1 contract (do not soften)

| Tool | Invoke | Output | Gate |
|------|--------|--------|------|
| Gitleaks | `gitleaks detect --no-git … --exit-code 1` | `gitleaks.json` | fail on leaks |
| Trivy | `trivy fs --scanners vuln --severity HIGH,CRITICAL --exit-code 1` | `trivy.json` | fail on HIGH/CRITICAL |
| Checkov | `checkov` CLI (**never** `python -m checkov`) | `checkov.json` | fail on failed checks |
| Semgrep | `semgrep scan --config p/ci --json` | `semgrep.json` | fail on findings |

Run all scanners with `continue-on-error: true`, open/dedupe issues (`issues: write`), then **Stage 1 gate** on the recorded exit codes.

## Actions hygiene

- Pin every `uses:` to a full 40-char commit SHA (`# vN` comment).
- Never start a bash heredoc body at column 0 inside `run: |` (breaks YAML). Prefer indented `python3 -c`.
- Avoid `curl | interpreter`; download to a temp file, then parse.

## Inference / autofix

Issues labeled `autofix:queued` are handled by skill `scanner-hotfix`: process first → laptop Ollama (auto-fallthrough) → habitat Qwen (mock until live) → DeepSeek as Frontier AI. One branch per issue into `dev`; after merge close/delete the issue and delete the branch. Include the PR token report table.

## Deploy

- Prod: human merge to `main`, then `frontier release-check` + deploy script, or `cloud-deploy.yml` after verify succeeds.
- Prefer `gcloud run services update --image` (or Cloud Build that updates the service) so Cloud SQL / secret mounts / sizing stay intact.

# Release process — feat → dev → main

Standard merge and deploy. Open-source tools do the scans. Frontier authorizes the laptop push and the production deploy. GitHub branch protection owns "no direct push to `dev` or `main`."

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
| Require PR + checks on `dev` and `main` | **GitHub** branch protection |
| Exam + ledger before a laptop push | Frontier `plan` → `apply` |
| Secrets / SAST / SCA / IaC / tests | GitHub Actions (`verify.yml`) using OSS CLIs |
| Mutate production | App `deploy.ps1` (or equivalent) |
| “This commit may go to production” | Frontier `release-check` inside that deploy |

Laptop deny-commit-on-dev/main is a hook catch, not a separate subsystem.

## App already using Frontier: new platform version before deploy

Frontier is a **tool plus copied templates**, not a git submodule of the app. When platform changes (feat → dev → main on `Wadek/frontier-platform`), the app does not `git pull` this repo into its tree. Prompt the **app** agent with this, then deploy only after it lands:

```text
frontier-platform has a new version. Do not deploy this app yet.

1. Work on feat/<slug> in THIS app. Never commit or push dev or main.
2. Install the new frontier binary (rebuild from Wadek/frontier-platform
   into $FRONTIER_RUNTIME/bin, or pull the release asset). Confirm:
   frontier version
3. From this app tree: frontier ready
   Apply every fail/advise hint (verify.yml, Gitleaks/Trivy/Semgrep/Checkov,
   deploy script, compose, GitHub protection on BOTH dev and main).
4. Diff this app's copied files against frontier-platform
   ship/templates/release/ and take the new process
   (feat -> PR -> dev -> PR -> main).
5. frontier plan && frontier apply && git push
   gh pr create --base dev --fill
6. Human merges to dev, then a second PR promotes dev to main.
   Deploy only from a checkout of main: frontier release-check, then
   the app deploy script.
```

`ready` still does not mutate the app. The agent or human applies the hints.

## Stages (no LLM)

| Stage | When | Tools |
|-------|------|--------|
| **0** | `git push` of `feat/*` or `dev` | Frontier `plan` → `apply` (OWASP v0 High/Critical, dirty tree, ledger) |
| **1** | PR into `dev` or `main` | Gitleaks, Semgrep, Trivy fs, Checkov, ESLint/Ruff |
| **2** | same PRs | Stack tests (`pytest -v`, `node --test`, `npm test`, `go test`) |
| **3** | PR | container build, Trivy image, Checkov on Dockerfile/compose |
| **4** | optional later | Playwright, ZAP, light k6 |
| **5** | after human merge to `main` | `deploy.ps1`: backup, up, health, rollback. `frontier release-check` first. |

Copy `ship/templates/release/` into the app on **first deploy** (hard-fail `verify.yml.example`, `scripts/ci/report_scanner_issues.py`, PR token-report template, optional `cloud-deploy.yml.example`). Copy `ship/templates/runner-docker/` to `ops/github-runner/` when using a Linux self-hosted runner. `frontier ready` lists what is still missing.

## CI/CD Runner Decision Matrix

Frontier-platform relies on GitHub Actions (`verify.yml` and deployment workflows). When configuring your workflow's `runs-on:` target, use the following decision matrix to determine whether to use GitHub-hosted runners or self-hosted runners.

### GitHub-Hosted Runners (`ubuntu-latest`)
**Default Choice.** Use this for 95% of standard cloud deployments.

*   **When to use:**
    *   Deploying to public cloud infrastructure (GCP Cloud Run, AWS, Vercel, Cloudflare) with accessible API endpoints.
    *   Running standard pipeline stages (SAST, unit tests, Docker builds).
    *   Pushing to cloud-native container registries (e.g., GCP Artifact Registry).
*   **Why:** Ephemeral (clean state every run), zero maintenance overhead, auto-scaling, and natively secure.
*   **Security:** Authenticate via Workload Identity Federation (OIDC) rather than pasting static service account JSON keys into GitHub Secrets.

### Self-Hosted Runners
**Specialized Choice.** Use when network, hardware, **or billing** requires it.

*   **When to use:**
    *   **Homelab / On-Premise:** Deploying to local infrastructure (e.g., a local homelab SQLite database, Proxmox cluster, or Raspberry Pi) that sits behind a strict NAT/firewall.
    *   **Heavy Compute:** Running ML/AI model training, GPU compilation, or massive memory jobs that exceed standard GitHub Actions quotas.
    *   **Strict VPC Access:** Connecting to internal databases inside a locked-down VPC without provisioning a VPN tunnel to GitHub's dynamic IP ranges.
    *   **Hosted minutes exhausted:** Private-repo GitHub-hosted minute caps / billing freezes. Keep the same workflow YAML; switch `runs-on` to self-hosted labels (e.g. `[self-hosted, Windows, X64]`). Do **not** leave `ubuntu-latest` jobs running in parallel.
*   **Security Warning (CRITICAL):** **Never** attach a self-hosted runner to a public repository without requiring manual approval for all outside contributors. Malicious pull requests can execute arbitrary code directly on your internal network.
*   **State Warning:** Runners are persistent. CI scripts must proactively tear down test containers, dangling volumes, and clear workspaces to prevent state-bleed between pipeline runs.
*   **Prefer Linux self-hosted when practical:** WSL2 Ubuntu (or a small Linux VM) with a user/system systemd runner avoids Windows PATH/`pwsh`/`checkov.cmd`/cp1252 issues. Install gitleaks/trivy into a stable bin dir on PATH; use `shell: bash` and `runs-on: [self-hosted, Linux, X64]`.
*   **Windows service pitfalls (if you must):** Prefer `shell: powershell` over `pwsh` under LocalSystem. Install via elevated `config.cmd --runasservice` / `RunnerService.exe`. Pin CLIs with absolute paths or `setup-python` Scripts on PATH.
*   **WSL keep-alive:** User systemd runners need the distro running (`wsl -d Ubuntu`). Enable lingering (`loginctl enable-linger`) with sudo once so user services survive logout; otherwise the runner goes offline when WSL idles.
*   **Scanner contract:** Gitleaks = secrets; Trivy fs = `--scanners vuln` when gitleaks already covers secrets; Checkov = `checkov` / `checkov.cmd` CLI (**never** `python -m checkov`). Stage 1 should write JSON reports and open/dedupe GitHub issues on the runner (`issues: write`) so humans (or a later **local** agent) plan fixes — do not burn cloud-agent tokens mid-CI.
*   **Autofix loop (local-first):** Issues labeled `autofix:queued` (especially gitleaks/trivy) are picked up by a local-first Frontier agent using skill `scanner-hotfix`: one branch per issue (`fix/vuln-<n>` / `fix/issue-<n>`), PR into `dev` with `hotfix` label and a **token report** (process / laptop Ollama / habitat Qwen mock / DeepSeek). After merge to `dev`, close+delete the issue when allowed and delete the feature branch. Inference cascade: process → laptop Ollama (auto-fallthrough) → habitat Qwen (mock until live) → DeepSeek as Frontier AI. Human still promotes `dev` → `main`.
*   **Pin Actions `uses:` to full commit SHAs** (with a trailing `# vN` comment). Mutable tags (`@v4`) trip Semgrep `github-actions-mutable-action-tag` and enable supply-chain retargeting.
*   **YAML `run: |` + bash heredoc:** Never start a heredoc body at column 0 inside a GitHub Actions `run: |` block — YAML ends the literal scalar early and the workflow fails before jobs run. Prefer indented `python3 -c '…'` (or keep every heredoc line indented with the block).
*   **Docker self-hosted runners:** Prefer Compose on Docker Desktop (no host sudo). Template: `ship/templates/runner-docker/` (pre-baked gitleaks/trivy/checkov/semgrep; HEALTHCHECK on `Runner.Listener`; documented CKV_DOCKER_8 skip). Look for `Runner.Listener` under `/actions-runner` as well as `~/actions-runner`.
*   **App Dockerfile:** non-root `USER` + `HEALTHCHECK` against `/health` (CKV_DOCKER_2 / CKV_DOCKER_3). Mount secrets under `/secrets/` and symlink — never onto `/app`.
*   **Copyable Stage 1:** `ship/templates/release/verify.yml.example` is the hard-fail contract (JSON reports → issue reporter → gate). Do not ship the old soft-fail / `pip install` every job shape.
*   **Post-merge cleanup:** `ship/templates/release/cleanup-merged.yml.example` closes/deletes linked issues and deletes the head branch after a PR merges to `dev`.
*   **Frontier AI:** `frontier ai route|complete|tokens` and `frontier enhance guard --call` use the local-first cascade (Ollama → habitat mock → DeepSeek). Spend is recorded in `.frontier/tokens.jsonl` for PR token reports (`frontier ai tokens`).
*   **DeepSeek key:** You create the key at DeepSeek; Frontier does not ship one. Resolution order everywhere the process runs: (1) process env `DEEPSEEK_API_KEY` — GitHub Actions secrets inject here on hosted **and** self-hosted runners when the workflow sets `env: DEEPSEEK_API_KEY: ${{ secrets.DEEPSEEK_API_KEY }}`; (2) local keystore via `frontier ai secrets set deepseek` (Windows DPAPI under `%LOCALAPPDATA%\frontier\secrets`, Unix `~/.config/frontier/secrets`). Check with `frontier ai secrets status` (never prints the value).
*   **`frontier ready`:** detects reporter script, Stage 1 gate, SHA-pinned Actions, PR token-report template, runner-docker copy, and cleanup-merged workflow.

## First deploy (detect, do not mutate)

From the app tree:

```text
frontier ready
```

JSON + text: pass / fail / advise per item, with a fix hint. Frontier does not write the app. An agent or human applies the hints, then re-runs `ready`.

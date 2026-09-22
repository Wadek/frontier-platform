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

Copy `ship/templates/release/` into the app on **first deploy**. `frontier ready` lists what is still missing.

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
*   **Autofix loop (local-first):** Issues labeled `autofix:queued` (especially gitleaks/trivy) are picked up by a **local** Frontier agent using skill `scanner-hotfix`: patch on `fix/vuln-<n>` / `fix/hotfix-<n>`, PR into `dev` with `hotfix` label, comment on blocked product PRs that they depend on that hotfix. Human still promotes `dev` → `main`.

## First deploy (detect, do not mutate)

From the app tree:

```text
frontier ready
```

JSON + text: pass / fail / advise per item, with a fix hint. Frontier does not write the app. An agent or human applies the hints, then re-runs `ready`.

# Runner efficiency checklist

Detail for `runner-efficiency`. Apply only the sections the measurement
points to.

## 1. Hardware and storage

- **CPU.** Prefer high single-core clock speed (consumer desktop parts) over
  many slower server cores (Threadripper/EPYC class). Compilation and
  single-threaded critical paths benefit most from clock speed.
- **Disk.** Use local NVMe SSDs. Git clone, package install, and container
  builds are disk-bound.
- **Memory.** Ensure enough RAM that heavy compiles and multi-stage Docker
  builds never swap. Swap use during a job is a finding.

## 2. Caching and environment

- **Docker layers.** Copy package manifests and install dependencies before
  copying source, so source edits do not invalidate the install layer. Keep a
  persistent local BuildKit cache so unchanged layers are reused.
- **Dependency directories.** Persist `node_modules`, `.cargo`, `vendor`, the
  Go module cache, and pip caches across runs, keyed strictly on the matching
  lockfile (`package-lock.json`, `pnpm-lock.yaml`, `Cargo.lock`, `go.sum`,
  `requirements*.txt`).
- **Pre-baked tooling.** Bake heavy SDKs and base images into the runner image
  or VM template instead of downloading them per job.

## 3. Concurrency and state

- **Contention.** Match the number of concurrent jobs to physical cores, RAM,
  and disk throughput. More parallel jobs than the host can serve slows all of
  them.
- **Single-runner queues.** One Windows service runner processes one job at a
  time. A long Stage 1 or a backlog makes later PRs look hung; check the
  Actions queue before blaming the service.
- **Ephemeral runners.** Configure runners to take one job and exit
  (`./config.sh --ephemeral` or a container scale set) so no state leaks
  between jobs. Keep caches on a volume so this does not cost speed.

## 4. Scanner tooling on the runner

- Pre-bake gitleaks, trivy, checkov, and semgrep into the runner image or a
  stable tools directory (for example Frontier `runtime/bin`) so Stage 1 does
  not `pip install` on every job.
- Prefer the Compose template in `ship/templates/runner-docker/` (Docker
  Desktop, Linux labels, no host sudo). Use separate venvs for checkov and
  semgrep; pin `setuptools==75.8.2` for semgrep's `pkg_resources`.
- After `actions/setup-python`, put `$pythonLocation\Scripts` on PATH, or call
  `checkov.cmd` / `semgrep.exe` by full path. Never assume a user-profile
  install exists for the LocalSystem account.
- Split responsibilities: gitleaks for secrets; `trivy fs --scanners vuln`
  for dependencies; checkov for IaC; semgrep for SAST. Persist JSON and open
  GitHub issues from the workflow. Triage with local models or humans, not
  cloud agents in the hot path (see `scanner-hotfix`).
- Copy `ship/templates/release/verify.yml.example` and
  `ship/templates/release/scripts/ci/report_scanner_issues.py` so the gate
  matches the baked scanners.

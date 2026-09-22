---
name: runner-efficiency
description: >-
  Use this skill to diagnose and optimize slow self-hosted CI/CD runners. Contains hardware, caching, and concurrency best practices.
user-invocable: true
---

# Self-Hosted Runner Efficiency

When tasked with optimizing or diagnosing slow self-hosted CI/CD runners, apply the following analysis steps and recommendations.

## 1. Hardware and Storage Assessment
*   **CPU Profiling:** Ensure the runner prioritizes high single-core clock speeds (e.g., overclocked consumer CPUs like Intel i9) over expensive multi-core server processors (e.g., Threadripper/EPYC). Compilation and single-thread critical paths benefit most from high clock speeds.
*   **Disk I/O:** Verify the runner uses fast local NVMe SSDs. Git cloning, package installation, and container generation are heavily disk-bound.
*   **Memory Swapping:** Check for adequate RAM. Heavy compilation or multi-stage Docker builds will thrash the disk and destroy performance if memory limits force swap usage.

## 2. Caching and Environment Optimization
*   **Docker Build Layers:** Analyze `Dockerfile`s to ensure they are structured efficiently (e.g., copying package definitions and installing before copying source code). Set up a persistent local Docker build cache so unchanged layers are instantly reused.
*   **Dependency Directories:** Ensure heavy directories (`node_modules`, `.cargo`, `vendor`) are mapped or persisted across job runs, and strictly keyed to their respective lockfiles (`package-lock.json`, `pnpm-lock.yaml`, `Cargo.lock`).
*   **Pre-baked Tooling:** Rather than downloading heavy SDKs or base images in every pipeline execution, advise baking these dependencies directly into the runner's base OS image or VM template.

## 3. Concurrency and Resource Tuning
*   **Contention Limits:** Diagnose if too many parallel jobs are configured for the runner. Tune the maximum concurrency limits to perfectly match the physical capacity of the machine to prevent CPU/IO starvation.
*   **Single-runner queues:** One Windows service runner processes **one job at a time**. A long Stage 1 or a backlog makes later PRs look “hung” — check the Actions queue before blaming the service.
*   **Ephemeral State:** To prevent state pollution and hidden caching bugs, recommend configuring runners to be **ephemeral**. They should handle a single job and self-destruct (using `--ephemeral` flags or containerized scale sets), guaranteeing a pristine clean state for the next build without sacrificing pre-baked tool speed.

## 4. Scanner tooling on the runner
*   **Pre-bake** gitleaks, trivy, checkov, and semgrep into the runner image or a stable tools dir (e.g. Frontier `runtime/bin`) so Stage 1 does not `pip install` on every job.
*   After `actions/setup-python`, put `$pythonLocation\Scripts` on PATH (or call `checkov.cmd` / `semgrep.exe` by full path). Never assume a user-profile install exists for LocalSystem.
*   Split responsibilities: gitleaks → secrets; trivy fs → `--scanners vuln`; checkov → IaC CLI; semgrep → SAST. Persist JSON and open GitHub issues from the workflow — triage with **local** models / humans, not cloud agents in the hot path.

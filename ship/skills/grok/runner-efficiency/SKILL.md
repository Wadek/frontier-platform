---
name: runner-efficiency
description: >-
  Diagnose and speed up slow self-hosted GitHub Actions runners: measure where
  job time goes, then fix hardware, caching, concurrency, runner state, and
  scanner tooling in that order of evidence. Use when self-hosted CI is slow,
  jobs look hung or queued, Stage 1 scanners reinstall on every run, or the
  user asks to tune a runner.
user-invocable: true
---

# Self-hosted runner efficiency

Scope: self-hosted GitHub Actions runners (Windows service, Linux, or the
Docker Compose template in `ship/templates/runner-docker/`). Not for
GitHub-hosted runners, where hardware is fixed.

## Pstack root

- `principle-explain-the-number`: measure step timings before proposing a
  fix, and explain the number after.
- `benchmark-checklist`: a before/after timing is only reported after it is
  vetted (same job, same commit, warm vs cold stated, repeated).
- `principle-make-operations-idempotent`: ephemeral runners and pre-baked
  tools give each job the same starting state.
- `principle-boundary-discipline`: caches keyed on lockfiles; state does not
  leak across jobs.

## Steps

1. **Measure first.** Pull per-step durations for the last few runs of the
   slow workflow and find the dominant steps. Check the queue before blaming
   the runner: a single runner takes one job at a time, so a backlog looks
   like a hang. Commands: [diagnosis commands](references/diagnosis.md).
2. **Classify the dominant cost** and apply the matching fix from the
   [checklist](references/checklist.md):

   | Symptom | Likely cause | Section |
   |---|---|---|
   | Jobs wait with status `queued` | Too few runners or one long job | Concurrency |
   | Install/setup steps dominate | Tools downloaded each run | Pre-baked tooling, scanners |
   | `docker build` dominates | Layer order or no persistent cache | Caching |
   | Dependency install dominates | Cache missing or not keyed to lockfile | Caching |
   | Compile/test dominates, CPU pegged | CPU or contention | Hardware, concurrency |
   | High iowait, swap in use | Disk or RAM | Hardware |
   | Flaky results between runs | State left by earlier jobs | Ephemeral runners |

3. **Change one thing**, re-run the same workflow on the same commit, and
   compare step durations.
4. **Report** before/after per step with the run URLs.

## Done when

- [ ] The report names the dominant step(s) with measured durations and run
      URLs, before any recommendation.
- [ ] Each applied change has a before/after timing on the same workflow and
      commit, cold or warm cache stated.
- [ ] Queue time and run time are reported separately.
- [ ] Scanner tools resolve from a pre-baked path (no per-job `pip install`)
      when Stage 1 scanners are in scope.

## Conflicts

- Ephemeral runners discard warm caches; keep caches in a mounted volume or a
  pre-baked image so pristine state does not cost speed.
- Hardware advice (single-core clock, NVMe) is a recommendation to the user;
  never change host hardware or service configuration without being asked.

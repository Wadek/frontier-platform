---
name: long-running-background-tasks
description: >-
  Required before you start, watch, wait on, or report on anything that keeps
  running after launch: background jobs, watchers, scheduled loops, CI, pull
  requests, training runs, dev servers, long builds. Done means the job
  finished and its outputs were verified, not that it was queued. Use when
  about to launch, supervise, inspect, diagnose, or report the state of such
  work; saying where a running job stands counts as working on it.
---

# Working with long-running background tasks

Submitting or queuing a job does not complete the task. Done means the job
reached a successful terminal state and you verified it produced every
expected output. Watch it through to that point.

## Pstack root

- `principle-prove-it-works`: "finished" is read from the job's real state and
  outputs, not inferred from silence or a log's mtime.
- `principle-outcome-oriented-execution`: the outcome is the verified result;
  relaunch after a fix instead of stopping at the first failure.
- `principle-never-block-on-the-human`: diagnose and relaunch reversible jobs
  without asking.
- `principle-make-operations-idempotent`: watchers survive restarts and
  re-pointing without duplicate wake-ups.

pstack has no monitoring playbook; this skill fills that gap. Tool names
below are generic: **monitor** is any background process whose stdout wakes
the agent (Grok `monitor`, a background shell the agent awaits). **loop** is
any scheduled agent run (Grok `scheduler_create`, a routine, a cron'd
subagent).

## Choose the watcher

| | monitor | loop |
|---|---|---|
| Event is | a predicate a script can evaluate: exit code, API field, port answering | a judgment from reading logs: loss unstable, garbage output, throughput sagging, faulty nodes |
| Cost | free while silent | one wake-up per cycle |
| Examples | PR CI red, merge conflict, PR mergeable, dev server up/down, job terminal | training health, data pipeline quality, deployment 4xx/5xx and latency |

Use both when a job has deterministic terminal states *and* qualitative
health (training runs).

## Monitor contract

Plan before writing:

1. Enumerate every terminal condition (PR merged or closed; job completed or
   failed unrecoverably).
2. Enumerate every non-terminal wake-up (CI failed, conflict, new comment of
   any of the four PR comment types, PR became mergeable, loss diverged,
   instance preempted).
3. Decide what the watcher must re-target (a new head SHA after a push, a new
   pid after a restart).

Then the script must:

- Poll on a fixed cadence: 30 s for GitHub APIs (rate limits), 5 s for a
  local port.
- Print to stdout **only** wake-up lines starting with `DONE:`, `FAILED:`,
  `CANCELLED:`, or `ACTION_REQUIRED:`. Exit on the first three.
- Fail early: emit `FAILED:` as soon as any required condition is terminally
  failed.
- Send all diagnostics, its own and the job's, to a log file keyed on the
  watched pid, job id, or PR, so no log line is mistaken for a token and two
  watchers never clobber each other.
- Live in a durable path (for example `~/.grok/long-running-background-tasks/`)
  with a descriptive name (`ci_watcher.sh`).

Reference implementation for PRs: [scripts/watch_pr.sh](scripts/watch_pr.sh).
Worked examples (dev server, PR, training run):
[references/examples.md](references/examples.md).

## Handling wake-ups

- Act on errors immediately. Read the logs and find the root cause before
  resubmitting; long jobs are expensive to rerun blind.
- After a fix: re-push (PR) or requeue (job), then re-point or restart the
  monitor and loop at the new target.
- With several jobs running, start debugging the first failure at once;
  minimize total wall time.
- Flaky CI: rerun once. A second identical failure is a real failure.
- PR comments: ask the user once how to handle each comment type (auto-fix,
  reply, resolve), store the answer for future sessions, then apply it.

## Restarts and cleanup

- Every relaunch gets a fresh monitor; kill the one it replaces.
- Monitors may have a hard runtime cap (Grok: 10 hours). Restart before the
  cap if the job is still running.
- On a terminal state, kill the monitor, the loop, and any dead job.
- An outdated monitor (attached to an old run) that wakes you is killed.
- Leave no detached workers, duplicate pollers, or orphaned watchers.

## Done when

- [ ] The job's terminal state was read from the source of truth (`gh`,
      `sacct`, process table, HTTP status), and the report quotes it.
- [ ] Every expected output exists and was checked (artifact present, metric
      at target, endpoint answering).
- [ ] No monitor, loop, or job process from this task is still running unless
      the user asked to keep it (`ps`/scheduler list checked).
- [ ] Each failure seen during the run has a stated root cause and the
      relaunch that followed.

## Conflicts

- `token-thrift`: watchers stay silent until actionable; that is thrift.
  Polling the job from the agent turn by turn is not.
- `scanner-hotfix` says not to spend agent turns watching scanner CI because
  the self-hosted runner re-checks on push. That wins for scanner PRs; use a
  silent monitor only if the user wants a wake-up.

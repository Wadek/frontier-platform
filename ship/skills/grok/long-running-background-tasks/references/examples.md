# Long-running task examples

Worked examples for `long-running-background-tasks`. Tool calls use Grok
names (`run_terminal_command`, `monitor`, `scheduler_create`); map them to
the host agent's equivalents.

## Example: maintaining a dev server

Suppose you just made some changes to the user's frontend code and would like to spin it up to test it. While you're testing it, you would like to know if the server goes down for any reason (if you made changes, for example). This is a great use case for the monitor tool. First start a React app in a background terminal, which reports its pid:

```
run_terminal_command(
  command="echo $$; npm start >~/.grok/long-running-background-tasks/dev_server_$$.log 2>&1",
  timeout=0,
  description="Start the dev server on :3000",
  background=true
) # prints the pid 123456789
```

There are no terminal states, but there are two actionable events from the watcher: `http://localhost:3000` coming up and going down.

You want to be notified as soon as it is up so you can use it, and you want to be woken when it stops so you can address it — fix the code, restart the server, free the port. Both down paths quote the tail of the server's log, so the wakeup arrives with the error already in hand. Polling every 5s is a good cadence since you want to quickly fix it if it is down, and it is not expensive to poll.

```bash
#!/usr/bin/env bash
# Log paths are keyed on the server pid so two dev servers cannot clobber each other.
PID=$1 URL=http://localhost:3000 DIR=~/.grok/long-running-background-tasks
LOG=$DIR/dev_server_$PID.log
up() { curl -fsS -o /dev/null --max-time 5 "$URL" 2>>"$DIR/watch_dev_server_$PID.log"; }
why() { tail -20 "$LOG" | tr -d '\r' | tr '\n' ' '; }

until up; do
  kill -0 "$PID" 2>/dev/null || { echo "FAILED: pid $PID exited before $URL came up: $(why)"; exit 1; }
  sleep 5
done
echo "ACTION_REQUIRED: $URL is up"

while :; do
  sleep 5
  up || { echo "ACTION_REQUIRED: $URL is down: $(why)"; until up; do sleep 5; done; }
done
```

To spawn the monitor, point it at the pid `npm start` reported:

```
monitor(
  command='bash ~/.grok/long-running-background-tasks/watch_dev_server.sh 123456789',
  description="Watch the dev server on :3000"
)
```

On the first wakeup the app is ready to use and you can start testing it. After that, you may be woken up if the server goes down. When you do:
- Read `~/.grok/long-running-background-tasks/dev_server_<pid>.log` for the full compile error or port conflict and fix it.
- If the dev server survived, hot reload picks the fix up and the watcher recovers on its own, so leave the monitor running.
- If the process is dead, restart it with a fresh `run_terminal_command`, then kill this monitor and spawn a new one against the new pid.

## Example: PR monitoring

PR monitoring should stop when the PR is merged or closed, and it should wake up on:
- For comments, you should ask users how they would like to resolve them (which ones to address automatically, whether to respond / resolve, etc.), then store the decision for new future sessions. There are four different types of PR comments.
  - Inline review threads
  - Review bodies
  - Conversation comments
  - Commit comments
- Failing CI - most can be fixed in the code but sometimes they may be flaky. For flaky CI, just rerun them.
- Merge conflicts
- The PR becoming ready to merge — CI passes, no merge conflicts, reviews complete, no blockers, no unresolved inline threads. That wakeup is `kind: mergeable`. Tell the user they can merge now. It is not terminal; keep watching until merge or close.

The script bundled with this skill is [`watch_pr.sh`](../scripts/watch_pr.sh). It needs an authenticated `gh` CLI (`gh auth login`). Default poll is 30 s; a lower rate may hit GitHub rate limits at higher concurrency. It covers conversation comments, review bodies, and inline review comments; commit comments are not watched.

```
monitor(
  command='bash ~/.grok/long-running-background-tasks/watch_pr.sh "https://github.com/<owner>/<repo>/pull/123"',
  description="Watch PR 123"
)
```

## Example: training monitoring

For monitoring long training jobs, the terminal conditions are when the job either fails or completes. Further, many issues cannot be caught deterministically, so use a subagent to run a script for most of the high-level update, then scan the logs for errors such as faulty nodes or CUDA warnings.

The prompt below assumes a `watch_training.py` helper that parses nanoGPT stdout out of a Slurm run directory. It is not shipped in this repository; substitute your own status script that prints one `STATUS` and one `ISSUES` line.

```md
Keep an eye on the nanoGPT run nanogpt_owt_124m_r1 (slurm job 1234567) and tell me what changed since last time. Worktree ~/worktrees/nanogpt-owt-124m.

1. python <path-to>/watch_training.py --run-name nanogpt_owt_124m_r1 --job-id 1234567 --target-val 3.00

How to report:
- Start with the job's state and if it changed, e.g. "still RUNNING" or "PENDING -> RUNNING".
- Copy the ISSUES line and the STATUS line from the monitor exactly as it is written, numbers and all. STATUS will include where the run is heading: how far along, iters left, the eta, and whether val loss is still coming down. We want val loss at or under 3.00.
- If the state didn't change, tell me that in one line and stop.
- If the job has finished or died, give me the state, the exit code, and the last ISSUES, and run these two for me:
  sacct -j 1234567 -D -X -P -o JobID,State,ExitCode,DerivedExitCode,Elapsed,NodeList,FailedNode
  tail -20 ~/train_logs/slurm-nanogpt_owt_124m_r1-1234567.out

Keep <=10 lines: state, progress and eta, loss trend, then anything that smells like training going unstable (loss spiking, NaN, val loss creeping back up) or the cluster misbehaving (a drained or failing node, OOM, NCCL errors, a preemption, iters suddenly slower than usual).
```

Launch the scheduler on a 10 min interval, and spawn a monitor to be notified when the job crashes or succeeds.

```
scheduler_create(
  interval="10m",
  prompt: "...", # Copy from the prompt above.
)
monitor(
  # Wakes on pending/running/finished changes after the first observation.
  # First poll is silent unless the job is already COMPLETED or failed —
  # PENDING matches *ING, so treating an empty previous state as a
  # transition would wake immediately on a just-submitted job.
  command='J=1234567; p=; while :; do s=$(sacct -nXP -j $J -o State|head -1); s=${s%% *}; s=${s:-PENDING}; if [ "$s" != "$p" ]; then case $s in COMPLETED) echo "DONE: $J"; exit;; PENDING|RUNNING|CONFIGURING|COMPLETING|REQUEUED|SUSPENDED|RESIZING) [ -n "$p" ] && echo "ACTION_REQUIRED: $J $p -> $s";; *) echo "FAILED: $J $s"; exit 1;; esac; fi; p=$s; sleep 30; done',
  description="Watch training run"
)
```

Launch both, so that we know when the job completes. When the monitor reports that the job has completed, cancel both the monitor and the scheduler, then determine next steps. Then:
- If it succeeded, simply present the results to the user.
- If it failed, debug the issue thoroughly, then fix and relaunch the job with new monitors attached. Make sure to have a thorough understanding of the issue before resubmitting, since training jobs are expensive. But always relaunch jobs so you don't waste the user's time.

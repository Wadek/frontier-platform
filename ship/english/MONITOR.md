# Monitor — audit agent behavior against the ship directives

`frontier monitor` is the supervision family of Frontier Ship. Whenever an
agent uses frontier-platform (shim, hooks, or the CLI), every consequential
action is sealed into the ledger (F0). The monitor replays that evidence and
verifies the agent followed the directives. It never rewrites history; it reads
the ledger and reports.

```
  English  →  this file + the directive list below      (policy)
  Go       →  internal/monitor/monitor.go               (what runs)
  Witness  →  internal/monitor/monitor_test.go          (replay rules, checked)
```

---

## Commands

```powershell
frontier monitor             # audit this repo's ledger
frontier monitor all         # audit every ledger under $FRONTIER_RUNTIME/ledgers
frontier monitor status      # recent monitor.* seals
frontier monitor directives  # print the D0-D7 reference set
```

Alias: `frontier audit` is the same command.

Every run seals one `monitor.audited` row (F0 evidence that the examination
happened). The monitor is **advise-only** by default, like Optimize. Set
`FRONTIER_MONITOR_BLOCK=1` to make `frontier monitor` exit non-zero when the
verdict is `violation` or `tampered`.

## Verdicts

| Verdict | Meaning |
|---------|---------|
| `clean` | every recorded action followed the directives |
| `watch` | denied attempts were recorded (the gate held, but an agent tried) |
| `violation` | at least one directive was broken (D1-D4, chaos inject) |
| `tampered` | the ledger hash chain does not verify (F0 breach) |

## The directives (v0)

| ID | Directive |
|----|-----------|
| D0 | F0 evidence: ledger hash chain is intact (prev_hash links, entry_hash recomputes) |
| D1 | No remote effect without a fresh sealed gate: `push.authorized` needs a fresh `gate.passed` that itself follows a fresh `plan.passed` (same branch+HEAD, 15 min TTL) |
| D2 | No ship from `main`/`master`: plan/gate/push/commit rows must be on feature branches |
| D3 | No soft bypass: `push.soft_allow` (`FRONTIER_SOFT=1`) is forbidden for real ship |
| D4 | No ship with untriaged High/Critical under V: `gate.passed` must not follow a blocking `exam.owasp` for the same branch+HEAD |
| D5 | Denied attempts are incidents: `push.deny`, `apply.deny`, `commit.deny_main`, `plan.failed`, `gate.failed` |
| D6 | Hygiene service reachable during inspect: `hygiene.service_down` is an advisory |
| D7 | Runtime chaos stays dry: chaos injection is not implemented; `chaos_denied` is an advisory |

D0 and the time-based freshness window (15 min) are verified in the Go runtime;
`internal/monitor/monitor_test.go` witnesses the state machine of D1-D5.

## Watching the ledgers (event-driven, not polling)

`scripts/frontier-monitor-watch.ps1` watches `$FRONTIER_RUNTIME/ledgers` with a
PowerShell `FileSystemWatcher` and runs `frontier monitor all` after each new
or changed `ledger.jsonl`. It prints one line per examined event and logs to
`$FRONTIER_RUNTIME/ledgers/monitor/watch.log`. Cheap watcher, one-shot work — the
same doctrine as the ship monitor skills. Install it as a logon scheduled task
only with operator confirmation:

```powershell
powershell -ExecutionPolicy Bypass -File ./scripts/frontier-monitor-watch.ps1 -Once
```

## Skills and agents

Bundled skills and agents ship inside frontier-platform and are available through the same CLI:

```powershell
frontier skills
frontier skills show dogfood
frontier agents
```

Trees: `skills/` and `agents/` (see their README files). Monitor directives D1–D5 are the machine-checkable subset of rules those skills teach (`dogfood`, `regular-git`, `test-as-you-go`, `local-first`).

## What the monitor is not

- Not a new push gate: it does not block `git push`; plan/apply already do that.
- Not a model: every check is a deterministic replay of sealed rows (no tokens).
- Not a rewrite tool: it reports findings; triage and remediation stay human or
  agent work, then re-run the monitor to confirm the fix.

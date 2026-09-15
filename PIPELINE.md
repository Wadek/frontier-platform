# Permanent local pre-GitHub build pipeline

This is the **only** ship process. It applies to every repo, every branch, every client. Nothing goes to GitHub until it passes.

Canonical code: `ship/` (history from frontier-ship)  
Runtime: `$FRONTIER_RUNTIME` (defaults to `./runtime` in this repo)  
Enforcement: set `core.hooksPath` to `$FRONTIER_RUNTIME/hooks` (or use lefthook in-repo)

---

## The process (do not skip, do not soften)

```text
1. Work on a feature branch   (never main / master)
2. Commit a clean tree
3. frontier hygiene           (optional: AI provenance on the changeset)
4. frontier plan              (OWASP Guard + Hygiene line + push rules; fail closed)
5. frontier apply             (seal gate.passed from a fresh plan)
6. git push                   (hook re-runs plan → apply, then allows the remote)
7. Open a PR into main        (human merge; do not push main)
```

Same commands:

```powershell
git checkout -b frontier/topic
git add -A
git commit -m "msg"
frontier hygiene   # advise; FRONTIER_HYGIENE_BLOCK=1 to fail closed
frontier plan      # must exit 0
frontier apply     # must exit 0
git push -u origin HEAD
gh pr create --fill
```

Helper (same thing, refuses main):

```powershell
powershell -File ./ship/scripts/dogfood-push.ps1
```

---

## What the gate actually checks today

| Check | Blocks ship? | Notes |
|-------|----------------|-------|
| Feature branch (not `main`/`master`) | **yes** | Commit on main is also refused |
| Clean working tree | **yes** | `.frontier/` dirt is ignored |
| OWASP v0 High/Critical | **yes** | Built-in `frontier guard` / ScanTree |
| Fresh `plan.passed` then `gate.passed` | **yes** | Ledger under `$FRONTIER_RUNTIME/ledgers/…`, 15 min TTL |
| Secret-surface / Checkov / enhance | advise | Not a hard block yet |
| Hygiene (H) | advise | `frontier hygiene`; block only if `FRONTIER_HYGIENE_BLOCK=1` |
| Slim (S) | no | Planned, not enforced |
| Layer A tests | **manual** | Required when a stack playbook exists; not in the binary gate |
| Layer B browser | **manual** | See `ship/english/O_VERIFY.md` |
| GitHub Actions `verify.yml` | remote | Does not replace the local gate |

`FRONTIER_SOFT=1` is **forbidden** for real ship. Hooks force `FRONTIER_SOFT=0`.

---

## What is *not* a legal path to GitHub

- `git push --no-verify`
- `git -c core.hooksPath= push`
- Calling a system `git.exe` to skip Frontier hooks
- GitHub MCP `push_files` / `create_or_update_file`
- Agent `merge_pull_request` onto main
- Direct commits or pushes to `main` / `master`

---

## Verify / Optimize (when the repo has tests)

After Guard, and **before** opening the PR, run the stack playbook verbosely and paste the log. Standard: `ship/english/O_VERIFY.md`.

---

## Where evidence lives

Ledger is **outside** the work tree (does not dirty the diff):

`$FRONTIER_RUNTIME/ledgers/<hash>/ledger.jsonl`

---

## Turn the gate off (you should not)

```powershell
git config --global --unset core.hooksPath
```

That is a policy break, not a convenience flag.

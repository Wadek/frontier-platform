# Permanent local pre-GitHub build pipeline

Repository rule, not a skill: it applies to every push from a repo that uses
the Frontier gate (this repo and repos initialized with `frontier`). Do not
invent a shorter path.

- **Spec:** [PIPELINE.md](../../../../PIPELINE.md) at the repo root.
- **Tool:** Frontier Ship (`frontier` / `git frontier`) on PATH or under
  `$FRONTIER_RUNTIME/bin`.
- **Helper:** `ship/scripts/dogfood-push.ps1` runs plan, apply, and push and
  refuses `dev`, `main`, and `master`.

## Pstack root

- `principle-prove-it-works`: the gate's exit codes and the sealed ledger row
  are the evidence, not "I ran it".
- `principle-encode-lessons-in-structure`: the sequence is enforced by hooks
  (`lefthook.yml` → `scripts/frontier-gate.ps1`), not by memory.

Condition: applies only when the repo has the Frontier gate (a
`lefthook.yml` or `core.hooksPath` pointing at Frontier hooks, or a
`PIPELINE.md`). Elsewhere, `regular-git` is the habit to follow.

## Required sequence before GitHub

1. Feature branch `feat/<slug>` (never commit or push `dev`, `main`, or
   `master`).
2. Clean commit; `git status --porcelain` is empty.
3. `frontier hygiene` checks the changeset for AI provenance (optional service
   at `http://127.0.0.1:8765`). Advise by default;
   `FRONTIER_HYGIENE_BLOCK=1` makes it blocking.
4. `frontier plan` exits 0 (OWASP Guard and push rules).
5. `frontier apply` exits 0 (seals `gate.passed`; plan must be fresh, 15 min
   TTL).
6. `git push` with hooks; they re-run plan and apply.
7. Open a PR into `dev`. A human merges. `dev` is promoted to `main` by a
   second PR.
8. Agents do not merge to `dev` or `main`.

`FRONTIER_SOFT=1` is forbidden for a real ship.

## Forbidden

- `git push --no-verify`
- Overriding `core.hooksPath` (including `git -c core.hooksPath= push`)
- Calling a system `git` in a way that skips Frontier hooks
- GitHub MCP `push_files`, `create_or_update_file`, or `merge_pull_request`
- Committing or pushing on `dev`, `main`, or `master`

## Tests

When the repo has a stack playbook or test suite, run it verbosely
(`pytest -v`, `go test ./... -v`) and paste the log. Standard:
[O_VERIFY.md](../../../english/O_VERIFY.md). Tests do not substitute for
`frontier plan` / `frontier apply`.

## Done when

- [ ] `frontier plan` and `frontier apply` exit codes (0) are in the report,
      from a run after the last commit.
- [ ] The push went through hooks from a feature branch; the PR targets `dev`.
- [ ] Verbose test log pasted when a suite exists.
- [ ] If the gate or a hook could not run (missing tool, unsupported OS), the
      report says so explicitly instead of claiming the gate passed.

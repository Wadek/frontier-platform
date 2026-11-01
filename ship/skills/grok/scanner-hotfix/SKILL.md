---
name: scanner-hotfix
description: >-
  Agent half of the scanner autofix loop: take a CI scanner issue (gitleaks,
  trivy, checkov, semgrep) labeled `scanner` + `autofix:queued`, fix it on one
  branch per issue with a test, open a PR into `dev` with a token report,
  notify blocked PRs, and close the loop after merge. Use when Stage 1
  scanners fail, when such issues are open, or when the user asks to clear
  scanner findings without burning cloud tokens.
user-invocable: true
---

# Scanner → issue → hotfix branch (local-first)

The runner opens deduplicated issues; this skill fixes them. Runner side:
`ship/templates/release/verify.yml.example` plus
`ship/templates/release/scripts/ci/report_scanner_issues.py`.

## Pstack root

- `principle-fix-root-causes`: remove the secret or bump the dependency; a
  suppression is a fix only with a written reason.
- `tdd`: lock the fix with a regression check when a cheap path exists.
- `blast-radius`: a dependency bump or base-image change can break callers
  outside the diff; prove the one fact it is safe because of.
- `principle-sequence-verifiable-units`: one finding, one branch, one PR.

Conditions: requires the runner-side workflow and the `scanner`,
`autofix:queued` labels. Model choice follows `local-first`.

## Control rules

1. **Humans own merge to `main`.** Hotfix PRs target `dev`. Agents merge into
   `dev` only if repo policy allows it.
2. **Inference** follows the `local-first` cascade: process → laptop Ollama →
   habitat Qwen (mocked) → DeepSeek. No other cloud model unless the user
   names one.
3. Never auto-merge secret redactions to `main`. Never rewrite git history
   unless the user explicitly orders a history purge.
4. The Frontier ship path applies: feature branch → `frontier plan` /
   `frontier apply` → push → PR into `dev`
   ([local build pipeline](../rules/local-build-pipeline.md)).
5. **One branch per scanner issue**: `fix/vuln-<n>-<slug>` or
   `fix/issue-<n>-<slug>`.

## Loop

1. **Intake.** `gh issue list --label scanner --label autofix:queued --state
   open`. Read tool, file, line, rule/CVE, and workflow URL from the body.
   gitleaks and trivy are autofix candidates; checkov and semgrep may need
   human design. Swap the label to `autofix:in-progress` and comment
   `autofix: in progress`.
2. **Patch.** Branch from fresh `origin/dev`. Smallest safe fix:

   | Tool | Fix | Allowed suppression |
   |---|---|---|
   | gitleaks | Remove the secret, load from env or Secret Manager, redact docs. Treat the value as leaked: tell the user to rotate it. | None for a real secret. Test fixtures only, via allowlist with reason. |
   | trivy | Bump the package or base image to the fixed version. | `.trivyignore` entry with a dated reason and an expiry. |
   | checkov | Fix the IaC setting. | Inline `checkov:skip=<ID>:<reason>`. |
   | semgrep | Fix the code. | `nosemgrep: <rule-id>` with a reason comment. |

   Add or adjust a test that locks the fix when a cheap path exists. Re-run
   the same scanner locally on the changed files and confirm the finding is
   gone. Run the repo test suite verbosely.
   Commit: `fix(sec): <tool> <id> (#<issue>)`.
3. **PR into `dev`.** Title `fix(sec): <summary> (#<issue>)`; labels
   `hotfix`, `scanner`, and the tool name; body per
   [PR body template](references/pr-body.md), including `Closes #<issue>` and
   the token report.
4. **Notify dependents.** Find open PRs whose Stage 1 failed on the same
   finding. Comment once on each: hotfix PR `#N` unblocks the scanner gate;
   rebase onto `dev` after it merges.
5. **Close the loop after merge to `dev`.**
   1. Confirm the issue closed (via `Closes #n` or manually), label
      `autofix:done`.
   2. Delete the head branch (`gh pr merge --delete-branch` or delete the ref).
   3. Delete the issue (`gh issue delete <n> --yes`) only when the user or
      repo policy has approved issue deletion; it is irreversible. Closed and
      labeled `autofix:done` is the default.

Do not spend agent turns watching CI; the self-hosted runner re-checks on
push. If you must wait, use `long-running-background-tasks`.

## Runner-side contract

Stage 1 must:

1. Run scanners with their real CLIs (never `python -m checkov`).
2. Write JSON reports.
3. Open deduplicated issues with `autofix:queued` for gitleaks and trivy.
4. Fail the job in a final gate step after all scanners ran.

## Done when

- [ ] Local re-run of the scanner on the changed files shows the finding gone
      (output in the PR body).
- [ ] Test suite passed verbosely; regression check added or its absence
      explained.
- [ ] PR targets `dev`, links the issue with `Closes #<n>`, carries labels and
      the token report.
- [ ] Every dependent PR found in step 4 has one comment.
- [ ] After merge: issue closed and labeled, head branch deleted. A leaked
      secret was reported to the user for rotation.

## Conflicts

- `tdd` vs "add tests when practical": `tdd` decides. A version bump usually
  has no cheap regression test; the scanner re-run is the check.
- Issue deletion is destructive; default to close plus label.

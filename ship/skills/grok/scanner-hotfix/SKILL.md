---
name: scanner-hotfix
description: >-
  When CI scanners (gitleaks, trivy, checkov, semgrep) open GitHub issues,
  spawn a local-first subagent to patch on a feature branch into `dev`, mark
  hotfix/vuln status, and notify blocked PRs of the dependency.
user-invocable: true
---

# Scanner → issue → hotfix branch (local-first)

Use when Stage 1 scanners fail and issues are labeled `scanner` + `autofix:queued`
(or the user asks to clear gitleaks/trivy/checkov findings without burning cloud tokens).

## Control rules

1. **Human owns merge to `main`.** Hotfix PRs target **`dev`**. Agents may merge into `dev` only if the repo policy allows it.
2. **Inference cascade (local-first):** (a) tools/process with zero LLM tokens, (b) laptop Ollama (auto-fallthrough when GPU/model is insufficient), (c) Mac mini / habitat Qwen-class (mock until that host exists), (d) **DeepSeek API** as Frontier AI — never a different cloud coding model unless the user names one.
3. Do **not** auto-merge secret redactions to `main`. Do **not** rewrite git history unless the user explicitly orders a history purge.
4. Frontier ship path still applies: feature branch → `frontier plan` / `apply` → push → PR into `dev`.
5. **One feature branch per scanner issue** (`fix/vuln-<n>-…` or `fix/issue-<n>-…`). After that PR merges to `dev`: close (and delete when the API allows) the issue, and delete the head branch. The finding is done; keep neither artifact.

## Loop

### 1. Intake

- Find open issues with labels `scanner` and `autofix:queued` (gitleaks/trivy preferred for autofix; checkov/semgrep may need human design).
- Read the issue body: tool, file, line, rule/CVE, workflow URL.
- Comment on the issue: `autofix: in progress` (local agent) and remove `autofix:queued` or add `autofix:in-progress`.

### 2. Patch on a dedicated branch

- Branch from current `origin/dev`: `fix/vuln-<issueNumber>-<short-slug>` or `fix/hotfix-<issueNumber>-<short-slug>`.
- Apply the smallest safe fix:
  - **gitleaks:** remove committed secrets; load from env/Secret Manager; redact docs; rotate if the value was real.
  - **trivy:** bump the vulnerable package / base image; add `.trivyignore` only with a dated reason.
  - **checkov / semgrep:** fix the IaC/code smell or add a narrow suppress with reason.
- Add/adjust tests that lock the fix when practical.
- Commit message: `fix(sec): <tool> <id> (#<issue>)`.
- Run the repo test suite locally (verbose). Then `frontier plan` / `apply` / push.

### 3. PR into `dev`

- Title: `fix(sec): <summary> (#<issue>)`
- Labels: `hotfix`, `scanner`, tool name (`gitleaks` / `trivy` / …).
- Body links the issue (`Closes #<issue>` / `Fixes #<issue>`).
- Include a **token report** table in the PR body:

  | Route | Tokens |
  |-------|--------|
  | Local tools / process | … |
  | Laptop Ollama | … (0 if fallthrough) |
  | Habitat Qwen (mock until live) | … |
  | Frontier AI (DeepSeek) | … |

- State clearly: vulnerability/secret hotfix; blocked product PRs should rebase after merge.

### 4. Notify dependents

- Find open PRs whose Stage 1 failed for the same finding (or that still contain the secret/CVE).
- Comment on each: this hotfix PR `#N` unblocks scanner gate; please rebase onto `dev` after merge.
- Comment on the failing workflow/PR that opened the issue if known.

### 5. Close the loop

- After the hotfix merges to `dev`:
  1. Confirm the scanner issue is **closed** (`Closes #n` on the PR, or close manually).
  2. **Delete the issue** when the repo token can (`gh issue delete <n> --yes`). If delete is forbidden, closed + labeled `autofix:done` is the fallback.
  3. **Delete the feature branch** (`gh pr merge --delete-branch` or delete the ref).
- Do **not** spend cloud-agent turns watching CI; the self-hosted runner re-checks on push.

## Runner-side contract (apps)

Stage 1 should:

1. Run scanners with correct CLIs (never `python -m checkov`).
2. Write JSON reports.
3. Open deduped issues with `autofix:queued` for gitleaks/trivy.
4. Fail the job via a final gate after all scanners have run.

This skill is the **agent** half; the workflow is the **runner** half.

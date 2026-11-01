---
name: regular-git
description: >-
  Keep work on git and GitHub as a habit: a private repo when none exists,
  feature-branch commits in the same session as the work, and a factual answer
  to "is this pushed?". Use when the user says "we need to make sure we are
  regularly using git and github", "create private repo on github and commit",
  or "is this pushed to github?". Does not replace the Frontier pre-push
  pipeline.
user-invocable: true
---

# regular-git

Use git and GitHub in the same session as the work. Do not leave repos
uncommitted or unpushed.

## Pstack root

- `principle-sequence-verifiable-units`: each commit is a unit that lands on
  its own; the commit sequence reads as the argument for the change.

Frontier-specific part: in a Frontier-gated repo, pushing goes through
[local build pipeline](../rules/local-build-pipeline.md).

## Steps

1. **No repo yet:** `git init`, then `gh repo create <owner>/<name> --private
   --source . --remote origin`. Private is the default; public only when the
   user says so.
2. **Branch:** work on `feat/<slug>` (or the repo's prefix). Never commit or
   push `main`, `master`, or a protected `dev`.
3. **Commit** each verified unit with a message that states the change, not
   the activity ("fix: reject empty token", not "updates").
4. **Push** the branch with hooks enabled.
5. **"Is this pushed?"** Answer from git, not memory:

   ```sh
   git fetch origin
   git status -sb                       # ahead/behind vs upstream
   git rev-parse HEAD @{u}              # same SHA means pushed
   git status --porcelain               # empty means nothing uncommitted
   ```

   Reply "yes" or "no", the branch name, and the short SHA. If no upstream is
   set, the answer is "no".

## Done when

- [ ] Work is committed on a feature branch; `git status --porcelain` is empty.
- [ ] `HEAD` equals `@{u}` on the remote branch, and the report names the
      branch and short SHA.
- [ ] New repos are private unless the user asked otherwise.

## Conflicts

- In-session tests do not replace `frontier plan` / `frontier apply` before
  GitHub.
- Never skip local git hooks (`--no-verify`, empty `core.hooksPath`). If a hook
  blocks, report it.

---
name: dogfood
description: >-
  Exercise a change the way its user would, record the evidence, then push the
  feature branch. Use when the user says "dogfood it", "dogfood it and push to
  github", or asks you to try the thing you just built before shipping. Not
  the same as committing: the change must be run as a user runs it.
user-invocable: true
---

# dogfood

Run the thing you just built the way a user would, then push.

## Pstack root

- `principle-prove-it-works`: verify the real artifact, not a proxy.
- `blast-radius`: its "how sure are you" ladder is the evidence bar used below.
- `show-me-your-work`: for long or unattended runs, keep the evidence in a
  decision log.

This skill adds the user's-eye pass and the push step. In this repository,
dogfood also means shipping through Frontier itself
([ship/english/DOGFOOD.md](../../../english/DOGFOOD.md)).

## Steps

1. **Name the user path.** One sentence: who uses the change and what they do
   (CLI command, HTTP request, UI flow, library call).
2. **Run it as that user.** Use the shipped entry point (built binary,
   installed CLI, running server, public API), not an internal function or a
   unit test. Use realistic input, including one bad input.
3. **Record evidence.** Command and output, response body, or screenshot. Rate
   it on the `blast-radius` ladder; the bar is level 4 (ran it) or level 5
   (reproduced in the running app). Levels 1 to 3 are not dogfood.
4. **Fix what you found** before pushing. Re-run step 2 after each fix.
5. **Push the feature branch.** In a Frontier-gated repo, follow
   [local build pipeline](../rules/local-build-pipeline.md):
   `frontier plan` and `frontier apply` exit 0, then `git push` with hooks.
   The helper `ship/scripts/dogfood-push.ps1` runs plan, apply, and push and
   refuses `dev`, `main`, and `master`.
6. Open a PR into `dev` when the repo uses a `dev` integration branch.

## Done when

- [ ] The user path is written down in one sentence.
- [ ] Evidence at ladder level 4 or 5 is in the report (command plus output,
      or screenshot), including the bad-input case.
- [ ] Problems found during dogfood are fixed and re-run, or listed as open.
- [ ] The feature branch is pushed and the report states the branch name and
      remote; in a Frontier-gated repo, plan and apply exited 0.

## Conflicts

- Push only when the user asked for it ("dogfood it and push"). "Dogfood it"
  alone means steps 1 to 4 and a report.
- Never skip local git hooks to get the push through. If the gate fails,
  fix the cause or report the failure.

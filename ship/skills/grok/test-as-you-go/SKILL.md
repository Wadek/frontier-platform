---
name: test-as-you-go
description: >-
  Run the project's tests in the same turn as each code change, then the full
  suite once before claiming done. Use when implementing, editing, refactoring,
  or generating application code; when the user asks for tests as you go, TDD,
  a produce-and-verify loop, a full test suite, a smoke test, or how to
  structure a test; or when they run /test-as-you-go [on|off|command].
when-to-use: >-
  producing code, writing or editing source files, implementing a feature,
  refactoring, TDD, run tests as you go, test as you produce, full suite of
  tests, smoke test, smoke tests, structure a test, at least we have some
  testing done, /test-as-you-go
argument-hint: "[on|off|test command]"
user-invocable: true
---

# Test as you go

Tests run in the same turn as the change. Do not batch edits and only then
discover they fail.

## Pstack root

- `principle-sequence-verifiable-units`: each change is one unit that ends in a
  checked state before the next unit starts. This skill fixes the check: the
  smallest covering test.
- `tdd`: governs *adding* tests. A new test is written only when a cheap,
  local test path exists; otherwise use the closest executable check.
- `principle-prove-it-works`: "done" is a passing run on the real code, not
  "it compiles".

Frontier-specific part: the `.grok/test-on-stop` Stop hook and the reminder
that tests do not replace `frontier plan` / `frontier apply`. Both apply only
in a Grok workspace and a Frontier-gated repo respectively.

## Loop

1. Make one coherent change (one behavior, not the whole feature).
2. Run the smallest test command that covers it.
3. On failure: paste the failing output, fix, rerun. Do not start the next
   behavior while red.
4. Before claiming done, run the project's full test command once and keep
   the summary line (pass/fail counts) for the final report.

## Resolve the test command

In this order; stop at the first hit:

1. First non-comment, non-blank line of `.grok/test-on-stop` in the workspace
   root.
2. Auto-detect from the tree:

   | Marker | Command |
   |---|---|
   | `package.json` with `scripts.test` | `pnpm test` / `yarn test` / `npm test`, chosen by lockfile |
   | `Cargo.toml` | `cargo test` |
   | `go.mod` | `go test ./...` |
   | `pyproject.toml`, `pytest.ini`, `setup.cfg`, or `tests/` with Python | `pytest -v --tb=short` |
   | `Makefile` with a `test` target | `make test` |

3. None found: say so in one line and skip. Do not invent a runner.

Prefer a targeted invocation when the change has an obvious test:
`pytest path -k name`, `cargo test name`, `go test ./pkg -run Name`,
`npx vitest run path`.

## `/test-as-you-go`

| Argument | Effect |
|---|---|
| none or `on` | Write `.grok/test-on-stop` in the workspace root with the resolved command. This arms the Stop hook: on turn end, after this session wrote files, rerun it and block finishing on failure. Confirm the path and command. |
| `off` | Delete `.grok/test-on-stop`. |
| anything else | Write the argument verbatim as the command. |

Never create `.grok/test-on-stop` in the user's home directory.

## Done when

- [ ] Every behavior change was followed by a covering test run in the same
      turn, or the report says why none applies.
- [ ] The full test command ran once after the last edit and passed; its
      summary line is in the final report.
- [ ] Any test skipped as impractical is named, with the closest check used
      instead (per `tdd`).
- [ ] In a Frontier-gated repo, `frontier plan` and `frontier apply` still ran
      before push (see [local build pipeline](../rules/local-build-pipeline.md)).

## Conflicts

- `tdd` says skip impractical tests; this skill says always run tests.
  Resolution: *run* the existing tests every change; *add* new tests only
  where `tdd` finds a cheap path.
- A full suite that takes longer than about 10 minutes: run the targeted test
  per change and the full suite once at the end, in the background (see
  `long-running-background-tasks`).
- In-session tests never replace the Frontier pre-push pipeline.

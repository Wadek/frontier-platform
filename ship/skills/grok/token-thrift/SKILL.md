---
name: token-thrift
description: >-
  Keep agent output short: tools and patches over narration, no filler, stop
  when the gate can pass. Use when the user says tokens are wasted, output is
  too long, "stop wasting tokens", or asks for terse answers; also when
  instructing a verbose peer agent.
user-invocable: true
---

# token-thrift

Length rules for agent output. The canonical policy text is
[ship/english/TOKEN_THRIFT.md](../../../english/TOKEN_THRIFT.md), printed on
every successful `frontier plan` and `frontier apply`.

## Pstack root

- `unslop`: prose quality (no AI vocabulary, no slogans, no false groups of
  three). This skill adds hard length limits on top.
- `principle-minimize-reader-load`: the reader's time is the budget.
- `principle-guard-the-context-window`: route bulk output to files or
  subagents; keep summaries in context.

## Rules

1. Prefer tools and patches over narration.
2. No slogans, metaphors, or invented section titles.
3. Between tool calls, at most one short status sentence.
4. When `frontier plan` / `frontier apply` can pass, stop expanding scope.
5. If a peer agent is verbose, restate these rules once, then continue the
   task. Do not argue about it.
6. Large output (logs, diffs, listings) goes to a file; quote only the lines
   that matter.

## Done when

- [ ] No paragraph of narration between consecutive tool calls.
- [ ] The final report holds results, evidence pointers, and open items only;
      no restated plan and no closing summary of what the user already saw.
- [ ] Scope stopped growing once the gate could pass.

## Conflicts

- `poteto-mode` and `swarm` favor multi-agent depth. When the user asked for
  fewer tokens, thrift wins: run single-agent unless told otherwise.
- Thrift never removes required evidence (test output, gate exit codes,
  links). Shorten the prose, not the proof.

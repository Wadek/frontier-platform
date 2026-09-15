---
name: token-thrift
description: >-
  Keep agent output short: prefer tools over prose, no filler, stop when the
  gate can pass. Use when the user says tokens are wasted, output is too long,
  or "stop wasting tokens".
user-invocable: true
---

# token-thrift

## Rules

1. Prefer tools and patches over narration.
2. Do not invent slogans, metaphors, or cute section titles.
3. Between tool calls, at most one short status sentence.
4. When `frontier plan` / `apply` can succeed, stop expanding scope.
5. If another agent is verbose, restate these rules once and continue the task.

See `ship/english/TOKEN_THRIFT.md`.

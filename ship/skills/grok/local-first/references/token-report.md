# Token report format

Every PR whose change was produced through the local-first cascade carries
this table in its body. Generate it with `frontier ai tokens` (reads
`.frontier/tokens.jsonl` in the repo) instead of typing numbers by hand.

| Route | Tokens |
|-------|--------|
| Local tools / process | 0 |
| Laptop Ollama | <n> (0 if fallthrough) |
| Habitat Qwen (mock until live) | 0 (mocked fallthrough) |
| Frontier AI (DeepSeek) | <n> |

Rules:

- Every row is present, even when it is 0. A missing row hides a tier that
  was skipped.
- Tier 2 reads `0 (mocked fallthrough)` until the habitat host is live, so the
  report maps future local capacity.
- `frontier ready` checks that the repo's PR template contains a token report
  (`pr_token_report`); the template ships in
  `ship/templates/release/pull_request_template.md`.

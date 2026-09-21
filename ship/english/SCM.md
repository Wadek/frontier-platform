# SCM onboarding — `frontier scm`

**Separate from Learn.** Customers who have no Git/GitHub need a place to put code before Learn/Guard/Slim/Optimize.

## Commands (intent)

| Command | Purpose |
|---------|---------|
| `frontier scm status` | Detect: git work tree? remotes? host (GitHub/GitLab/none)? |
| `frontier scm init` | Local only: `git init`, ignore suggestions (human confirms) |
| `frontier scm connect` | Guide: create remote / set `origin` (human auth — no silent cloud create) |

## Order

```
  frontier scm status|init|connect   # if needed
  frontier learn …
  frontier guard …
```

For dogfood on Wadek GitHub: `scm status` should report green (already managed).

## Rules

- Never create a remote repository without human authentication.
- Seal ledger events (`scm.detected`, `scm.initialized`) for local actions only.
- Learn classify assumes a project folder; SCM decides whether it is under version control yet.

## Status

**Phase A:** doctrine + CLI stub.

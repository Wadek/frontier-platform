# Dogfood — take our own medicine

We develop **Frontier Ship** using Frontier itself.

That is how we test the tool: if we cannot ship our own changes through `plan → apply → push`, customers will not either.

## Rules for this repo

1. **No casual push to `dev` or `main`.** Work on `feat/…` branches.  
2. **`FRONTIER_SOFT=0`** for real ship attempts (soft is for learning only).  
3. Always:

```text
git frontier learn      # optional: classify before change
git frontier guard      # optional visibility (security)
git frontier hygiene    # optional: AI provenance on the changeset
git frontier plan       # must exit 0
git frontier apply      # must exit 0
git push -u origin HEAD
```

4. Open a PR into `dev` (human merge). Promote `dev` to `main` with a second PR. Merging those PRs may use GitHub UI / `gh pr merge` so we do not invent a special "commit on protected branches" bypass in the tool.  
5. If plan fails on our own tree (OWASP V), **fix or triage** — do not soft-allow to save time.

## Why this matters

| Skip dogfood | Dogfood |
|--------------|---------|
| Tool bitrots | Failures surface on us first |
| Docs lie about the flow | Docs match reality |
| Customers hit unknown pain | We feel the same gate |

## Helper

```powershell
powershell -File scripts\dogfood-push.ps1
# runs plan → apply → push on current branch (refuses dev/main/master)
```

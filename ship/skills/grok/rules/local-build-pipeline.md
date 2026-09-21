# Permanent local pre-GitHub build pipeline

Every session, every repo, every push. Do not invent a shorter path.

**Spec:** `./PIPELINE.md` (repo root) or `ship/../PIPELINE.md` in the monorepo  
**Tool:** Frontier Ship (`frontier` / `git frontier`) on PATH or under `$FRONTIER_RUNTIME/bin`

## Required sequence before GitHub

1. Feature branch (never commit or push `dev` / `main` / `master`).
2. Clean commit.
3. `frontier hygiene` inspects the changeset for AI provenance (optional service at `http://127.0.0.1:8765`). Advise by default.
4. `frontier plan` must exit 0 (OWASP Guard + push rules).
5. `frontier apply` must exit 0 (seals `gate.passed`).
6. `git push` — hooks re-run plan/apply. Never `--no-verify`.
7. Open a PR into `dev`. Human merges. Promote `dev` to `main` with a second PR.
8. Agents do not merge to `dev` or `main`.

`FRONTIER_SOFT=1` is forbidden for real ship.

## Forbidden

- `git push --no-verify`
- Overriding `core.hooksPath`
- Calling system `git` in a way that skips Frontier hooks
- GitHub MCP `push_files`, `create_or_update_file`, or `merge_pull_request`
- Pushing or committing on `dev`/`main`/`master`

## Tests

When the repo has a stack playbook or test suite, run it **verbose** locally (`pytest -v`, etc.) and paste the log. Doctrine: `english/O_VERIFY.md`. Tests are not a substitute for `frontier plan`/`apply`.

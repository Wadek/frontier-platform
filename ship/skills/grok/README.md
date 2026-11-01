# Grok skills

Agent skills shipped with Frontier Platform. Each skill builds on a
[pstack](https://github.com/cursor/plugins/tree/main/pstack) skill or principle (its *pstack
root*) and adds Frontier-specific practice on top. The root supplies the
general discipline; the skill supplies the trigger, the concrete steps, and
the Frontier conditions.

Every `SKILL.md` has the same shape: frontmatter (`name` equal to the folder,
a `description` that says when to use it), a `## Pstack root` section, the
rules or steps, a `## Done when` checklist, and a `## Conflicts` section. Long
detail lives in `references/`; executable helpers in `scripts/`. The
structure is enforced by `ship/internal/repo/skills_test.go`, which runs in
`go test ./...` and CI.

## Index

| Skill | What it does | Pstack root | Frontier-specific condition |
|---|---|---|---|
| [test-as-you-go](test-as-you-go/SKILL.md) | Run the covering test after each change; full suite before done | `principle-sequence-verifiable-units`, `tdd`, `principle-prove-it-works` | `.grok/test-on-stop` Stop hook (Grok workspaces) |
| [dogfood](dogfood/SKILL.md) | Run the change as its user would, record evidence, push | `principle-prove-it-works`, `blast-radius`, `show-me-your-work` | Push through the Frontier gate |
| [token-thrift](token-thrift/SKILL.md) | Hard output-length rules | `unslop`, `principle-minimize-reader-load`, `principle-guard-the-context-window` | Policy text in `ship/english/TOKEN_THRIFT.md` |
| [long-running-background-tasks](long-running-background-tasks/SKILL.md) | Monitor or loop until verified done; PR watcher script | `principle-prove-it-works`, `principle-outcome-oriented-execution`, `principle-never-block-on-the-human`, `principle-make-operations-idempotent` | Tool names (`monitor`, `scheduler_create`) are Grok's; map to the host |
| [local-first](local-first/SKILL.md) | Tools, then local model, then cloud fallback; token report | `principle-laziness-protocol`, `principle-build-the-lever`, `principle-guard-the-context-window` | Ollama, habitat (mocked), DeepSeek are parameters; `frontier ai` implements the cascade |
| [scanner-hotfix](scanner-hotfix/SKILL.md) | Scanner issue → one fix branch → PR into `dev` → close the loop | `principle-fix-root-causes`, `tdd`, `blast-radius`, `principle-sequence-verifiable-units` | Needs the Stage 1 scanner workflow and `autofix:queued` labels |
| [runner-efficiency](runner-efficiency/SKILL.md) | Measure, then fix slow self-hosted runners | `principle-explain-the-number`, `benchmark-checklist`, `principle-make-operations-idempotent`, `principle-boundary-discipline` | Self-hosted GitHub Actions runners only |
| [regular-git](regular-git/SKILL.md) | Commit and push as a habit; answer "is it pushed?" from git | `principle-sequence-verifiable-units` | Push via the Frontier gate where present |
| [cloudflare](cloudflare/SKILL.md) | Load Cloudflare WAF guidance only for confirmed Cloudflare users | `principle-attack-the-premise` | Overrides the installed Cloudflare plugin skill |

Repository rule (not a skill): [rules/local-build-pipeline.md](rules/local-build-pipeline.md),
the fixed pre-GitHub sequence for Frontier-gated repos (roots:
`principle-prove-it-works`, `principle-encode-lessons-in-structure`).

The shared `learning` process skill lives with the compositional skills in
`control/templates/compositional_skills/learning/`.

## Conflicts between skills

| When both apply | Resolution |
|---|---|
| `tdd` (skip impractical tests) and `test-as-you-go` (always test) | Run existing tests after every change; add new tests only where `tdd` finds a cheap path. |
| `poteto-mode` / `swarm` (multi-agent depth) and `token-thrift` | If the user asked for fewer tokens, thrift wins; run single-agent. |
| `long-running-background-tasks` (watch to done) and `scanner-hotfix` (do not watch scanner CI) | Scanner PRs: the self-hosted runner re-checks; no agent watching unless the user asks. |
| `dogfood` ("push") and a request that only says "dogfood it" | Push only when asked. |
| Any skill and the Frontier gate | The gate is never bypassed. If it cannot run, say so. |
| pstack `arena` / `architect` (multi-model panels) and a single-model host | Run them single-model unless told otherwise. |

## Adding a skill

1. Create `<name>/SKILL.md` with the shape above; `name` must equal the
   folder.
2. Name at least one pstack root in `## Pstack root`.
3. Add a row to the index table.
4. Run `go test ./ship/internal/repo/ -run Skills -v`.

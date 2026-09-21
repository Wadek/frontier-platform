# Architecture

English policy is `ship/english/`. These diagrams are the same process, generic (no host names).

| File | What |
|------|------|
| [01-context.puml](01-context.puml) | Laptop, GitHub, deploy host |
| [02-branch-flow.puml](02-branch-flow.puml) | `feat → dev → main` |
| [03-sequence-push.puml](03-sequence-push.puml) | `plan` → `apply` → `git push` |
| [04-sequence-prod-release.puml](04-sequence-prod-release.puml) | Human merge, then deploy |
| [05-stage-tools.puml](05-stage-tools.puml) | OSS tools per stage |
| [06-deployment.puml](06-deployment.puml) | `release-check` then compose |

ASCII one-pager: `ship/english/ARCHITECTURE_ASCII.md`.

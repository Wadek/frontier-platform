# Scanner hotfix PR body

```markdown
fix(sec): <tool> <rule-or-CVE> in <path> (#<issue>)

Vulnerability/secret hotfix. Blocked product PRs should rebase onto `dev`
after this merges.

Closes #<issue>

## Finding
- Tool: <gitleaks|trivy|checkov|semgrep>
- Rule / CVE: <id>
- Location: <file>:<line>
- Workflow run: <url>

## Fix
<one or two sentences: what changed and why it removes the finding>
<if suppressed: the reason, and the expiry date for .trivyignore>

## Evidence
- Scanner re-run on changed files: <command> → <0 findings>
- Tests: <command> → <summary line>
- Secret rotation needed: <yes, owner notified | no>

## Token report
<paste `frontier ai tokens` output; format in
ship/skills/grok/local-first/references/token-report.md>
```

## Dependent PR comment

```markdown
Scanner gate on this PR fails on <tool> <id>. Hotfix #<N> fixes it in `dev`.
After #<N> merges, rebase this branch onto `dev`.
```

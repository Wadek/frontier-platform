# LEARNING — the .agent_learning format and taxonomy layout

Every session that touches a project ends with a debrief. No debrief, no `closed` (T4).

## `.agent_learning.jsonl` (append-only, ledger-style)

One JSON record per session, one record per line, appended in time order. Fields:

```json
{
  "session":   "YYYY-MM-DD-<slug>",
  "project":   "<project-dir-name>",
  "agent":     "<agent-id>",
  "provider":  "local|flash|pro|dsh",
  "ts":        "2026-09-14T08:30:00Z",
  "tldr":      "one paragraph",
  "learned":   ["fact…"],
  "skills_proposed": ["skill-slug"],
  "mappings_updated": ["path → meaning"],
  "attribution": {"repo": "…", "commit": "…", "session": "…"},
  "next":      ["…"]
}
```

Rules: append-only (never rewrite history — F0); no secrets (egress rules apply); control
renders it for humans (`control debrief show`). The teacher pass (`control data generate`)
condenses debriefs into knowledge docs and skills off-peak.

## Taxonomy branch per project (`.agent_<project>/`)

```
.agent_<project>/
    taxonomy.yaml                  # manifest: version, created_by, domain, providers,
                                   #   tier (untrusted→saint), lineage parent, expiry
    knowledge/
        architecture.md            # how the project is built
        mappings.md                # path → what-it-is fast lookup
        brief.md                   # one-page agent brief ("who am I, what do I own")
        # each doc carries document_outline + document: {repo, commit, patterns}
    compositional_skills/
        grounded/<skill>/SKILL.md  # context-bearing (references project docs)
        <skill>/SKILL.md           # ungrounded generic how-to
    sessions/.agent_learning.jsonl
    inbox/                         # CloudEvent-shaped transfer cards
```

## SKILL.md frontmatter (InstructLab field names)

```yaml
---
version: 1
name: <skill-slug>
task_description: |
  Prescriptive description of what the skill does — teacher-model grade detail.
created_by: <agent-id>
grounded: true|false
seed_examples:                    # min 1; ideally 3+ worked examples
  - context: |                    # grounded skills only; copy-paste from the project docs
      …
    question: …
    answer: …
attribution: {repo: …, commit: …, session: …}
---
# body — the how-to (progressive disclosure; SKILL.md conventions)
```

Templates: `control/templates/taxonomy.yaml`, `control/templates/SKILL.md`.

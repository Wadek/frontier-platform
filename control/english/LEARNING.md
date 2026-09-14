# LEARNING — the .agent_learning format and taxonomy layout

Every session that touches a plane ends with a debrief. No debrief, no `filed` (T4).

## `.agent_learning.jsonl` (append-only, ledger-style)

One JSON record per session, one record per line, appended in time order. Fields:

```json
{
  "session":   "YYYY-MM-DD-<slug>",          // unique session id
  "plane":     "<plane-dir-name>",           // which project this session touched
  "pilot":     "<agent-id>",                 // who did the work (e.g. waka-cli session id)
  "runway":    "local|flash|pro|dsh",        // which runway flew
  "ts":        "2026-09-14T08:30:00Z",       // RFC3339 UTC
  "tldr":      "one paragraph",
  "learned":   ["fact…"],                    // what this session proved or discovered
  "skills_proposed": ["skill-slug"],         // candidates for the teacher pass
  "mappings_updated": ["path → meaning"],    // plane-map deltas
  "attribution": {"repo": "…", "commit": "…", "session": "…"},
  "next":      ["…"]                         // handoff / follow-up items
}
```

Rules: append-only (never rewrite history — F0); no secrets (egress rules apply); the tower
renders it for humans (`tower debrief show`). The teacher pass (`tower data generate`)
condenses debriefs into knowledge docs and skills off-peak.

## Taxonomy branch per plane (`.agent_<plane>\` at the wakalabs root)

```
.agent_<plane>\
    taxonomy.yaml                  # manifest: version, created_by, domain, runways,
                                   #   tier (untrusted→saint), lineage parent, expiry
    knowledge\
        architecture.md            # how the plane is built
        mappings.md                # path → what-it-is fast lookup
        brief.md                   # one-page pilot brief ("who am I, what do I own")
        # each doc carries document_outline + document: {repo, commit, patterns}
    compositional_skills\
        grounded\<skill>\SKILL.md  # context-bearing (references plane docs)
        <skill>\SKILL.md           # ungrounded generic how-to
    sessions\.agent_learning.jsonl
    inbox\                         # CloudEvent-shaped handoff cards
```

## SKILL.md frontmatter (InstructLab field names)

```yaml
---
version: 1
name: <skill-slug>
task_description: |
  Prescriptive description of what the skill does — teacher-model grade detail.
created_by: <pilot-id>
grounded: true|false
seed_examples:                    # min 1; ideally 3+ worked examples
  - context: |                    # grounded skills only; copy-paste from the plane docs
      …
    question: …
    answer: …
attribution: {repo: …, commit: …, session: …}
---
# body — the how-to (progressive disclosure; SKILL.md conventions)
```

Templates: `control/templates/taxonomy.yaml`, `control/templates/SKILL.md`.

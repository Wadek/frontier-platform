# LEARNING — ledger, methods, and how other skills employ this

Every session that touches a project ends with a debrief. No debrief, no `closed` (T4 / C4).
Learning is a **platform process**, not a personal project skill. Other skills employ it;
they keep domain facts. This file owns method + ledger.

## Methods (what “excellent” means here)

Pick one primary method. Combine at most two. Domain skill supplies content.

| ID | Method | Use when | What the agent does |
|---|---|---|---|
| `debrief` | Session debrief | every completed run | append one `.agent_learning.jsonl` record |
| `generate` | Pretest / generation | before reading a source | ask 3 questions first; then read; score misses |
| `retrieval` | Active recall | facts, claims, APIs | hide answers; one prompt at a time; correct then next |
| `spaced` | Spaced review | anything that must persist | schedule miss → 1d / 3d / 7d in `next` |
| `interleave` | Interleaving | similar procedures | mix 2–3 related items; do not block by type |
| `elaborate` | Self-explanation | mechanisms, “why” | force “because…” tied to a project doc line |
| `faded_example` | Worked-example fading | procedures | full example → partial → solo |
| `deliberate` | Deliberate practice | a known miss | only the weak item; tight feedback |
| `teachback` | Teach-back | confirm mastery | learner restates in own words; agent scores gaps |
| `dual` | Dual source | video, diagram, code+doc | same idea from two encodings, then recall |
| `listen` | Listening drill | speech / YouTube | one caption segment; hide text; compare |
| `seed` | Seed harvest | teacher pass | turn debriefs into `seed_examples` on a SKILL.md |

Evidence bar: testing effect, spacing, interleaving, worked examples, generation,
deliberate practice. Do not add a method that is only re-reading or highlighting.

### YouTube is a source, not a skill

`listen` + `dual` may take a video URL. Fetch captions, then drill or digest.
Do not download media. Do not invent a per-user YouTube skill folder.

Fetch order: local caption script if present → published transcript search →
user pastes YouTube → More → Show transcript.

## Other skills employ this

During `running` or when teaching:

1. Domain skill hydrates its own taxonomy (facts, mappings, grounded skills).
2. If the task is learn / quiz / listen / teach / digest-a-source, select method IDs
   from the table and follow `templates/compositional_skills/learning/SKILL.md`.
3. On `completed` → write `debrief` before `closed`.
4. Teacher pass (`control data generate`) may run `seed` off-peak.

Grounded skills reference project docs in `seed_examples.context`.
Ungrounded how-tos live under `compositional_skills/<skill>/SKILL.md`.
The learning skill itself is ungrounded and reusable across projects.

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
  "methods":   ["debrief", "retrieval"],
  "skills_proposed": ["skill-slug"],
  "mappings_updated": ["path → meaning"],
  "attribution": {"repo": "…", "commit": "…", "session": "…"},
  "next":      ["spaced:miss-x@3d"]
}
```

Rules: append-only (never rewrite history — F0); no secrets (egress rules apply); control
renders it for humans (`control debrief show`). `methods` is optional for old records;
if present, every ID must be in the table above. `debrief` is implied by the write itself.

The teacher pass (`control data generate`) condenses debriefs into knowledge docs and skills
off-peak.

## Taxonomy branch per project (`.agent_<project>/`)

```
.agent_<project>/
    taxonomy.yaml                  # manifest: version, created_by, domain, providers,
                                   #   tier (untrusted→saint), lineage parent, expiry
                                   #   employs: [learning]
    knowledge/
        architecture.md
        mappings.md
        brief.md
    compositional_skills/
        grounded/<skill>/SKILL.md
        learning/SKILL.md          # copy from control/templates/compositional_skills/learning
        <skill>/SKILL.md
    sessions/.agent_learning.jsonl
    inbox/
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
employs: [learning]               # omit only if the skill never teaches or debriefs
seed_examples:
  - context: |
      …
    question: …
    answer: …
attribution: {repo: …, commit: …, session: …}
---
```

Templates: `control/templates/taxonomy.yaml`, `control/templates/SKILL.md`,
`control/templates/compositional_skills/learning/`.

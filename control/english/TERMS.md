# TERMS — the fleet glossary (standard vocabulary adopted)

Sources: [InstructLab taxonomy](https://docs.instructlab.ai/taxonomy/) ·
[InstructLab skills](https://docs.instructlab.ai/taxonomy/skills/file_structure/) ·
[InstructLab knowledge](https://docs.instructlab.ai/taxonomy/knowledge/file_structure/) ·
[RHDH Orchestrator](https://docs.redhat.com/en/documentation/red_hat_developer_hub/1.8/html/orchestrator_in_red_hat_developer_hub/assembly-orchestrator-rhdh)

| Standard term | In this fleet | Notes |
|---|---|---|
| taxonomy | the knowledge+skills tree; the fleet is the root, each plane is a **branch** | `.agent_<plane>\` |
| knowledge | facts about the plane: architecture, mappings, brief | `knowledge\*.md`, each pins `document: {repo, commit, patterns}` |
| compositional skills | performative how-tos | `compositional_skills\<skill>\SKILL.md` |
| grounded / ungrounded | skill references plane docs (context) or not | grounded lives under `compositional_skills\grounded\` |
| leaf node | one unit of knowledge or skill | one dir with its `SKILL.md`/doc |
| qna.yaml | skill/knowledge schema (`version`, `task_description`, `created_by`, `seed_examples`) | we use the same field names in SKILL.md frontmatter |
| seed_examples | worked examples that teach the pattern (`context` + `question`/`answer`) | debriefs distill into these |
| attribution.txt | provenance of a knowledge doc | our debriefs carry `attribution: {repo, commit, session}` |
| teacher model → synthetic data generation → student model | off-peak pro/max condenses sessions into taxonomy; local 7B/flash pilots consume it | `tower data generate` is the teacher pass |
| orchestrator | the tower itself (frontier-control) | |
| workflow / workflow instances | flight plan definition / one run | `plans\*.workflow.yaml` + journal rows |
| data index | the registry of planes, pilots, runs | `tower roster` |
| job service | scheduled one-shot tasks | `tower ops jobs` (Task Scheduler, never polling) |
| review | human approval gate before a run's destructive/expensive steps | `tower workflow review` |
| CloudEvents | the handoff card envelope | `{specversion, type, source, id, data}` in `inbox\` |
| ilab `taxonomy init/diff`, `data generate` | `tower taxonomy init|diff`, `tower data generate` | same verbs |

**Deliberately NOT adopted:** the fine-tuning pipeline (`ilab model train`). An 8 GB card is
not a training cluster; our students are pilots consuming the taxonomy, not fine-tuned models.
We take the vocabulary, the file shapes, and the teacher/student economics — not the training loop.

Fleet-only terms (no standard exists): **plane** (a project directory), **pilot** (its agent),
**runway** (a harness), **weather** (pricing windows), **speciation** (governed narrowing of
a pilot's scope — from the SSC kit).

# TERMS — platform glossary (standard vocabulary)

Sources: [InstructLab taxonomy](https://docs.instructlab.ai/taxonomy/) ·
[InstructLab skills](https://docs.instructlab.ai/taxonomy/skills/file_structure/) ·
[InstructLab knowledge](https://docs.instructlab.ai/taxonomy/knowledge/file_structure/) ·
[RHDH Orchestrator](https://docs.redhat.com/en/documentation/red_hat_developer_hub/1.8/html/orchestrator_in_red_hat_developer_hub/assembly-orchestrator-rhdh)

| Standard term | In this platform | Notes |
|---|---|---|
| taxonomy | the knowledge+skills tree; the platform is the root, each project is a **branch** | `.agent_<project>/` |
| knowledge | facts about the project: architecture, mappings, brief | `knowledge/*.md`, each pins `document: {repo, commit, patterns}` |
| compositional skills | performative how-tos | `compositional_skills/<skill>/SKILL.md` |
| grounded / ungrounded | skill references project docs (context) or not | grounded lives under `compositional_skills/grounded/` |
| leaf node | one unit of knowledge or skill | one dir with its `SKILL.md`/doc |
| qna.yaml | skill/knowledge schema (`version`, `task_description`, `created_by`, `seed_examples`) | same field names in SKILL.md frontmatter |
| seed_examples | worked examples that teach the pattern (`context` + `question`/`answer`) | debriefs distill into these |
| attribution.txt | provenance of a knowledge doc | debriefs carry `attribution: {repo, commit, session}` |
| teacher model → synthetic data generation → student model | off-peak pro/max condenses sessions into taxonomy; local/flash agents consume it | `control data generate` is the teacher pass |
| orchestrator | frontier-control itself | |
| workflow / workflow instances | workflow definition / one run | `plans/*.workflow.yaml` + journal rows |
| data index | the registry of projects, agents, runs | `control roster` |
| job service | scheduled one-shot tasks | `control ops jobs` (OS scheduler, never polling) |
| review | human approval gate before a run's destructive/expensive steps | `control workflow review` |
| CloudEvents | the transfer card envelope | `{specversion, type, source, id, data}` in `inbox/` |
| ilab `taxonomy init/diff`, `data generate` | `control taxonomy init\|diff`, `control data generate` | same verbs |

**Deliberately NOT adopted:** the fine-tuning pipeline (`ilab model train`). Consumer hosts
are not a training cluster; students are agents consuming the taxonomy, not fine-tuned models.
We take the vocabulary, the file shapes, and the teacher/student economics — not the training loop.

Platform-only terms (no standard exists): **project** (a working directory), **agent** (its
worker), **provider** (a model/harness adapter), **pricing** (peak/off-peak windows),
**child-agent forking** (human-gated narrowing of an agent's scope).

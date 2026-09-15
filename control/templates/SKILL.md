---
version: 1
name: <skill-slug>
task_description: |
  Prescriptive description of what this skill does — enough detail that a
  teacher model could generate new examples from it (InstructLab standard).
created_by: <agent-id>
grounded: true            # false when the skill never references project docs
seed_examples:
  - context: |            # grounded skills only: copy-paste from the project docs
      …
    question: …
    answer: …
attribution:
  repo: …
  commit: …
  session: …
---
# <skill-slug>

When to use, exact steps, and the boundary of what NOT to do.

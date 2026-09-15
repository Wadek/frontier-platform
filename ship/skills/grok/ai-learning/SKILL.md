---
name: ai-learning
description: >-
  Persist and resume project lessons in a workspace `ai_learning` directory so any
  AI can pick up where the last one left off. Use when the user says "point any ai
  at ai_learning", "load ai_learning", "generate the lesson for a future ai", or
  "put it in the workspace root as ai_learning". Not the Grok /learn skill.
user-invocable: true
---

# ai-learning

Write and load the drop-in lesson at `./ai_learning`. Any later AI should open that path and continue.

## Steps

1. Generate or update the lesson in the workspace root as `ai_learning`.
2. When asked to load it, read `./ai_learning` and continue from what is already there.
3. Do not run the Grok `/learn` skill for this. That skill changes harness skills from traces. It is not for summarizing sessions into notes or memory.

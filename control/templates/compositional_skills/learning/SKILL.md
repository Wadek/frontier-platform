---
version: 1
name: learning
task_description: |
  Shared learning process employed by other skills. Select a method ID from
  english/LEARNING.md, run source → digest → drill → review, and finish with a
  debrief record. Domain facts stay in the calling skill. YouTube and other
  media are sources for listen/dual, not standalone skills.
created_by: frontier-control
grounded: false
employs: []
seed_examples:
  - context: |
      Calling skill: ultra-training. User: "quiz me on Idaho 50 handling."
    question: Which method and what is the first output?
    answer: retrieval. One question, answer hidden. Facts from the training skill.
  - context: |
      User pastes a YouTube URL and says "dictation."
    question: Which method and what must not happen?
    answer: listen. One caption segment, text hidden. Do not download the video. Do not create a YouTube skill folder.
  - context: |
      Workflow state is completed.
    question: Can the run close?
    answer: No. Append a debrief record first (C4).
attribution:
  repo: Wadek/frontier-platform
  commit: dev
  session: 2026-09-25-learning-methods
---

# learning

Ungrounded process skill. Other compositional skills list `employs: [learning]`.

## Loop

1. Identify source (calling-skill docs, transcript, user text).
2. Select method IDs (`control/internal/learning.Select` or the table in `english/LEARNING.md`).
3. Run one drill unit. Hide answers first.
4. Record misses in `next`.
5. On session end: append `.agent_learning.jsonl` (`debrief`).

Stop after the step the human asked for.

## Boundaries

- Do not store domain facts here.
- Do not skip debrief to reach `closed`.
- Do not download media.
- Do not rewrite the learning ledger.

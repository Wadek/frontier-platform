# Frontier â€” what we have (ASCII)

Snapshot of the system as built. Simple English labels.

```
                         HUMANS / AI HOSTS
                    (read English Â· call tools)
                               |
         +---------------------+---------------------+
         |                                           |
         v                                           v
  +------------------+                    +----------------------+
  | english/         |                    | teach/               |
  | INIT LANGUAGE    |                    | corporate trial      |
  | AXIOMS MINIMALITY|                    | curriculum prompts   |
  | CONTRIBUTING     |                    +----------------------+
  +--------+---------+
           | meaning
           v
  +------------------+
  | haskell/         |   COMPUTE / PROOF (pure)
  | Frontier.Role    |
  | Frontier.Gate    |
  | Frontier.Laws    |
  +--------+---------+
           | same truth
           v
  +----------------------------------------------------------+
  | GO RUNTIME                                               |
  |                                                          |
  |  cmd/frontier-git  ===== named =====>  git.exe on PATH   |
  |       |                                  (the interface) |
  |       | passthrough most verbs                           |
  |       | guard: commit-on-main, push                      |
  |       v                                                  |
  |  internal/policy  role  ledger  gitx  egress             |
  |                                                          |
  |  cmd/frontier      ===== CLI =====>  humans and agents   |
  +------------+------------------------------+--------------+
               |                              |
               | FRONTIER_GIT_BIN             | commands:
               v                              |  learn guard hygiene
  +------------------------+                  |  plan apply push
  | Real Git (vendor)      |                  |  ledger demo explain
  | Git for Windows / etc  |                  +----------^-----------+
  +------------------------+                             |
               |                                         |
               v                                         |
        your work tree <------------------ FRONTIER_REPO |
               |                                         |
               |  evidence (default OUTSIDE tree)          |
               v                                         |
  +------------------------+                             |
  | $FRONTIER_RUNTIME/ledgers\   |  append-only JSONL          |
  |   <repo-id>/ledgerâ€¦    |  gate.passed / push.*       |
  +------------------------+                             |
                                                         |
  STUDY ONLY (not required to run)                       |
  $FRONTIER_RUNTIME/src/git   <--- upstream git/git clone ------+


  -------------------- CONTROL FLOW (push) --------------------

    git status / diff     â‰ˆ  Observer / Analyst
    git commit            â‰ˆ  Operator   (deny on main/master)
    git frontier gate     â‰ˆ  seal Clean(C) + feature branch
    git push              â‰ˆ  Executor   (needs fresh gate.passed)

    git frontier demo     â‰ˆ  SEE branch/dirty/gate/ledger/ladder


  -------------------- LANGUAGE RULE --------------------

    Draft:     any language (*)
    Admit:     English + Haskell + Go agree (F5)
    Prefer:    least code that still proves the result


  -------------------- NOT IN SCOPE (on purpose) --------------------

    custom OS / linux-ai kernel fork
    random fourth official runtime language
```

## One-screen postcard

```
                 ENGLISH (meaning)
                      |
                 HASKELL (proof)
                      |
        +-------------v-------------+
        |     GO RUNTIME            |
        |  git shim  |  frontier CLI|
        +------+-----+------+-------+
               |            |
          real git      agents / shell
               |
          work tree
               |
          ledger (disk)

   demo â†’ gate â†’ push
   F0 evidence Â· F1â€“F4 laws Â· F5 English+Haskell+Go
   minimality: least code that proves the result
```

# Frontier — what we have (ASCII)

Snapshot of the system as built. Simple English labels.

```
                         HUMANS / AI HOSTS
                    (read English · call tools)
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
           | meaning (policy)
           v
  +----------------------------------------------------------+
  | GO RUNTIME — implementation + witness                    |
  |                                                          |
  |  cmd/frontier-git  ===== named =====>  git.exe on PATH   |
  |       |                                  (the interface) |
  |       | passthrough most verbs                           |
  |       | guard: commit-on-main, push                      |
  |       v                                                  |
  |  internal/policy  role  ledger  gitx  egress             |
  |  *_test.go beside each — the machine-checked witness     |
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
  |   <repo-id>/ledger…    |  gate.passed / push.*       |
  +------------------------+                             |
                                                         |
  STUDY ONLY (not required to run)                       |
  $FRONTIER_RUNTIME/src/git   <--- upstream git/git clone ------+


  -------------------- CONTROL FLOW (push) --------------------

    git status / diff     ≈  Observer / Analyst
    git commit            ≈  Operator   (deny on main/master)
    git frontier gate     ≈  seal Clean(C) + feature branch
    git push              ≈  Executor   (needs fresh gate.passed)

    git frontier demo     ≈  SEE branch/dirty/gate/ledger/ladder


  -------------------- LANGUAGE RULE --------------------

    Draft:     any language (*)
    Admit:     English (policy) + Go test (witness) agree (F5)
    Prefer:    least code that still proves the result


  -------------------- NOT IN SCOPE (on purpose) --------------------

    custom OS / linux-ai kernel fork
    random second official runtime language
```

## One-screen postcard

```
                 ENGLISH (meaning / policy)
                      |
                 GO  (implementation
                      |  + _test.go witness)
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

   demo → gate → push
   F0 evidence · F1–F4 laws · F5 English+Go
   minimality: least code that proves the result
```

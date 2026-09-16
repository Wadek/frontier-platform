# Languages of Frontier

## For humans

Everything a person must understand is **simple English** under `english/`.  
English is the **policy** — what we mean.

## For the machine

**Go** builds the local `git` shim and `frontier` CLI. It is the only
implementation language. Boring on purpose.

## The witness

English alone drifts. Every policy rule carries a Go `_test.go` beside the code
that enforces it — table-driven, run in CI. Those tests are the **witness**: the
machine-checked statement that Go still means what English says.

## For contributors

**Any language is welcome at the edge.**  
Pushing into Frontier still requires reverse-engineering:

```
  any-language patch
        → English policy line
        → Go code + `_test.go` witness
```

See [MINIMALITY.md](MINIMALITY.md).

```
  English  →  what we mean   (policy)
  Go       →  what runs      (implementation + witness)
  *        →  what you may draft in — then reduce
```

We limit **volume**, not **human entry languages**.

---

## Why one language?

Multiplicative maintenance: every extra layer multiplies the places a rule can
drift. The old three-layer story (English + Haskell + Go) kept a Haskell proof
shelf that CI never ran; it is gone. Two layers remain:

| Layer | Job |
|-------|-----|
| **English** | what we mean — read by humans |
| **Go + `_test.go`** | what runs, and what is enforced — checked by CI |

## Why not C for the runtime?

Git’s engine is C. That does **not** mean Frontier’s policy layer should be C.

| | **Go** (keep) | **C** |
|--|----------------|--------|
| Job | `git` shim, `frontier` CLI, ledger I/O | Upstream `git/git`, kernels, tiny hooks |
| Safety | Memory-safe by default | Easy to violate F1 via bugs |
| Size of *our* code | Small for JSON/CLI/strings | Same features ⇒ more lines, more review |
| Witness story | The tests run wherever the code runs | C is not a better witness language than Go |
| Reviewer tax | Lower | Higher (Linus-patience people hate vibe-C) |

**Rule of thumb**

- **Go** — what runs beside git today, and the tests that witness it  
- **C** — only when we must touch git’s own code or OS guts later  

Rewriting the shim in C would add risk and volume without making the axioms clearer.  
If we ever patch upstream git, that patch may be C — and must still reverse-engineer to English + a passing Go test.

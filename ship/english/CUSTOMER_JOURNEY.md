# First customer journey — managing vibe code with Frontier

## Two different questions

| Question | Answer |
|----------|--------|
| **How do people contribute *to* Frontier?** | Draft anywhere → reverse-engineer to English + Haskell proof → Go only if runtime must change → keep it small ([CONTRIBUTING.md](CONTRIBUTING.md), [MINIMALITY.md](MINIMALITY.md)) |
| **How do people *use* Frontier on a codebase?** | Install `git` shim → point it at the project → gate/exam before push → grow `V` as you learn → fix High/Critical under `V` before ship |

Contributors improve the **tool**.  
Customers use the tool to discipline **their** repo (including CEO vibe code).

---

## Persona: CEO vibe-coded an app

They have something that “works in the demo” and is now a liability.

**Frontier does not magically rewrite the product.**  
It makes **shipping** obey rules: evidence, roles, examination against a supplied vuln set `V`, least change that clears the gate.

### Process (practical)

```
  0. FREEZE chaos
     - One owner
     - No more push-to-main from ChatGPT paste

  1. INSTALL Frontier git interface
     - git frontier explain
     - FRONTIER_SOFT=1 for week one (learn), then 0

  2. IMPORT the codebase (this lab = stand-in customer)
     - Fork or clone to a place YOU control
     - Never “fix” upstream intentional-vuln remotes in place of learning

  3. BASELINE exam (F4)
     - git checkout -b frontier/baseline
     - git frontier demo
     - git frontier gate
     - Read OWASP findings under current V
     - Ledger seals exam.owasp + gate.*

  4. TRIAGE (human + Analyst)
     - Map findings → real risk for THIS product
     - Accept that V is incomplete (pattern v0 ≠ full AppSec)
     - Write “must fix before any prod push” vs “later”

  5. MINIMAL remediation (Operator)
     - Smallest patches that clear High/Critical under V
     - Or shrink scope (delete dead vibe features)
     - Re-run gate until gate.passed

  6. SHIP only through Executor path
     - git push to YOUR remote / PR
     - main stays protected

  7. GROW V as you learn
     - New failure modes → English policy line → Haskell rule → Go check
     - Old Clean(C,V) invalid when V grows (re-exam)
```

### What “safer and more efficient” means here

| Safer | More efficient |
|-------|----------------|
| Can’t push main casually | Smallest fix that passes gate |
| Secrets / injection sinks blocked by V | Less token burn: one tool/step |
| Evidence of what was examined | Delete vibe features instead of polishing them |
| Soft → hard mode as habit forms | One ladder for humans and AIs |

Efficiency is **not** “AI writes faster.”  
It is **less code shipped, with a seal**.

---

## First customer lab

**Stand-in customer:** Damn Vulnerable Node Application (DVNA) — intentionally vulnerable Node/Express app for training.

- Upstream: `https://github.com/appsecco/dvna`  
- We work on a **local clone / our fork**, not “fixing” their teaching upstream as if it were our product  
- Treat it exactly like a CEO’s vibe app that suddenly needs rules

Success for this engagement:

1. Baseline gate shows OWASP V findings (or cleanly explains misses)  
2. Verbose axiom log shows F0/F4 (and F1 on main)  
3. A minimal “customer” change path is documented  
4. English note: what Frontier caught vs what still needs deeper AppSec

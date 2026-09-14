# First customer baseline — DVNA (stand-in for CEO vibe code)

**Customer fixture:** [appsecco/dvna](https://github.com/appsecco/dvna) (Damn Vulnerable NodeJS Application)  
**Local path:** `D:\frontier\projects\dvna-customer`  
**Branch:** `frontier/customer-baseline`  
**Date:** 2026-08-18

## Story

Treat DVNA as a product that “works in the demo” and is full of intentional (or vibe-shaped) hazards.  
Frontier’s job is not to become a full pentest suite on day one. It is to **stop careless shipping** until High/Critical under current `V` are examined.

## What we ran

```text
git frontier demo
git checkout -b frontier/customer-baseline
git frontier gate    # FRONTIER_VERBOSE axioms on
git frontier ledger
# no push to upstream appsecco/dvna
```

## Axioms observed

| When | Axiom | What happened |
|------|-------|----------------|
| demo on `master` | F1 boundary | Gate would fail: refuse ship to main/master |
| gate start | F0 | Ledger opened before any remote mutate |
| exam | F4 | OWASP Top 10 `V` scan at least once / at maximum for this tool |
| findings | F4 | High matches → block |
| seal | F0 | `exam.owasp` + `gate.failed` |
| continuity | F3 | No bypass — gate stays failed |

## Findings under V (after ignoring `*.min.js` noise)

Real app code (representative; from first scan of `core/appHandler.js`):

| Severity | OWASP | Signal |
|----------|-------|--------|
| High | A03 | SQL string concat from `req.body.login` |
| High | A03 | `exec('ping …' + req.body.address)` |
| High | A03 | `mathjs.eval(req.body.eqn)` |

Minified vendor JS initially produced false positives → **minimality fix:** skip `*.min.js` in the scanner (customer taught us).

## Customer outcome

- **Gate failed** — correct for a first “real rules” day.  
- CEO/vibe workflow cannot `git push` this tree until High/Critical under `V` are fixed or formally triaged (future).  
- Next customer steps: minimal patches **or** delete unused dangerous routes; re-gate; push only to **their** remote.

## Contribute vs use (reminder)

| Role | They do |
|------|---------|
| Frontier contributor | Improve English/Haskell/Go policy (e.g. better A03 without hitting jquery) |
| Customer team | Use `git frontier gate` on **their** fork until Clean(C,V) |

## Efficiency lesson

We did **not** ask anyone to rewrite DVNA.  
We made shipping **expensive** until the smallest clear hazards under `V` are faced — that is how Frontier manages vibe code volume.

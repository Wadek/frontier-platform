# Architecture (short)

See English docs for the real explanation.

```
  Human reads:     english/
  AI/compute:      haskell/src/Frontier/*.hs
  Machine runs:    Go cmd + internal
  Daily UX:        git  (Go shim → real git)
  Evidence:        ledger on disk (outside work tree by default)
```

Do not add another language without updating `english/LANGUAGE.md` and F5.

# Architecture (short)

See English docs for the real explanation.

```
  Human reads:     english/          (policy — what we mean)
  Machine checks:  Go *_test.go      (witness — what is enforced)
  Machine runs:    Go cmd + internal
  Daily UX:        git  (Go shim → real git)
  Evidence:        ledger on disk (outside work tree by default)
```

Do not add another language without updating `english/LANGUAGE.md` and F5.

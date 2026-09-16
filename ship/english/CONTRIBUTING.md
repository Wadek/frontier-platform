# Contributing

You do not need to be Linus.  
You do need to respect people who have Linus-level patience for noise.

## 0. Dogfood

Changes to **this** repo use Frontier on ourselves: `plan → apply → push` on a `feat/…` branch — not soft-bypass to `main`.  
See [DOGFOOD.md](DOGFOOD.md).

## 1. Make it small

Least code that still proves the result.  
See [MINIMALITY.md](MINIMALITY.md).

## 2. Draft in any language

Python, Rust, C, notes, sketches — fine for exploration.

## 3. Before you ask to merge / push

Reverse-engineer your change:

| Step | Output |
|------|--------|
| Claim | One English sentence in the PR / commit body |
| Policy | English line under `english/` **or** “no new law” |
| Witness | Go `_test.go` for the check, plus a paste of `git frontier demo` / `gate` / `ledger` |
| Size | Rough LOC; if large, say why |

## 4. Do not

- Add a second “official” implementation language without rewriting LANGUAGE.md  
- Dump generated trees, vendor blobs, or copy-pasted framework apps  
- Expand scope “while we’re here”

## 5. How to see your test

```text
git frontier demo
git frontier gate
git frontier ledger
```

If we cannot see it in those witnesses, it is not done.

# Contributing

You do not need to be Linus.  
You do need to respect people who have Linus-level patience for noise.

## 0. Dogfood

Changes to **this** repo use Frontier on ourselves: `plan → apply → push` on a `frontier/…` branch — not soft-bypass to `main`.  
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
| Proof | Haskell update under `haskell/` **or** “no new law” |
| Witness | Paste of `git frontier demo` / `gate` / `ledger` / `cabal test` |
| Size | Rough LOC; if large, say why |

## 4. Do not

- Add a fourth “official” runtime language without rewriting LANGUAGE.md  
- Dump generated trees, vendor blobs, or copy-pasted framework apps  
- Expand scope “while we’re here”

## 5. How to see your test

```text
git frontier demo
git frontier gate
git frontier ledger
```

If we cannot see it in those witnesses, it is not done.
